package postgres

import (
	"context"
	pg "github.com/LeoUraltsev/auth-service/internal/app/postgres"
	"github.com/LeoUraltsev/auth-service/internal/domain/users"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type UsersStorage struct {
	db  *pg.Postgres
	log *slog.Logger
}

type User struct {
	id           string
	name         string
	email        string
	passwordHash []byte
	isActive     bool
	createdAt    time.Time
	updatedAt    time.Time
}

func NewUsersStorage(db *pg.Postgres, log *slog.Logger) *UsersStorage {
	return &UsersStorage{db: db, log: log}
}

// Save метод для добавления новых пользователей, обновление существующих, проставить пользователю статус удаленного
// возвращает nil в случае успешного добавления
func (u *UsersStorage) Save(ctx context.Context, user *users.User) error {
	log := u.log
	log.Info("saving user to postgres")
	us := mapperToStorage(user)

	query := `INSERT INTO users (id, name, email, password_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id) DO UPDATE 
		SET name = EXCLUDED.name, 
		    email = EXCLUDED.email, 
		    password_hash = EXCLUDED.password_hash, 
		    is_active = EXCLUDED.is_active, 
		    updated_at = EXCLUDED.updated_at;
		`
	log.Debug("query to save user", slog.String("query", query))

	_, err := u.db.Pool.Exec(ctx, query, &us.id, &us.name, &us.email, &us.passwordHash, &us.isActive, &us.createdAt, &us.updatedAt)
	if err != nil {
		//todo: доп проверка на ошибку уникальности
		log.Error("failed to save user to db ", slog.String("error", err.Error()))
		return err
	}
	log.Info("user saved successfully", slog.String("id", us.id))
	return nil
}

//todo: пагинация, возможно сортировка

func (u *UsersStorage) Get(ctx context.Context, id uuid.UUID) (*users.User, error) {
	log := u.log
	query := `SELECT id, name, email, password_hash, is_active, created_at, updated_at FROM users WHERE id = $1;`
	var user User
	err := u.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.id,
		&user.name,
		&user.email,
		&user.passwordHash,
		&user.isActive,
		&user.createdAt,
		&user.updatedAt,
	)
	if err != nil {
		log.Error("failed to get user by id ", slog.String("id", id.String()))
		return nil, err
	}

	dUser, err := mapperToDomain(user)
	if err != nil {
		log.Error("failed to convert user to domain", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("successful getting user by id ", slog.String("id", id.String()))
	return dUser, nil
}

func (u *UsersStorage) GetAll(ctx context.Context) ([]*users.User, error) {
	log := u.log
	query := `SELECT id, name, email, password_hash, is_active, created_at, updated_at FROM users;`
	rows, err := u.db.Pool.Query(ctx, query)
	if err != nil {
		log.Error("failed to get all users", slog.String("error", err.Error()))
		return nil, err
	}
	usersList := make([]User, 0)
	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.id,
			&user.name,
			&user.email,
			&user.passwordHash,
			&user.isActive,
			&user.createdAt,
			&user.updatedAt,
		)
		if err != nil {
			log.Error("failed to get all users", slog.String("error", err.Error()))
			return nil, err
		}
		usersList = append(usersList, user)
	}
	dUsers := make([]*users.User, 0, len(usersList))
	for _, user := range usersList {
		mMser, err := mapperToDomain(user)
		if err != nil {
			log.Error("failed to convert user to domain", slog.String("error", err.Error()))
			continue
		}
		dUsers = append(dUsers, mMser)
	}
	log.Debug("len result", slog.Int("len", len(dUsers)))
	return dUsers, nil
}

func (u *UsersStorage) ExistsByEmail(ctx context.Context, email users.Email) (bool, error) {
	log := u.log
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_active = TRUE);`
	log.Debug("query to check if user exists", slog.String("query", query))
	var exists bool
	err := u.db.Pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		log.Error("failed to check if user exists by email", slog.String("email", email.String()))
		return false, err
	}
	log.Debug("exists user", slog.Bool("res", exists))
	return exists, nil
}

func mapperToStorage(u *users.User) User {
	return User{
		id:           u.ID().String(),
		name:         u.Name().String(),
		email:        u.Email().String(),
		passwordHash: u.Password().Hash(),
		isActive:     u.IsActive(),
		createdAt:    u.CreatedAt(),
		updatedAt:    u.UpdatedAt(),
	}
}

func mapperToDomain(u User) (*users.User, error) {
	id := uuid.MustParse(u.id)
	name, err := users.NewName(u.name)
	if err != nil {
		return nil, err
	}
	email, err := users.NewEmail(u.email)
	if err != nil {
		return nil, err
	}
	passwordHash, err := users.NewPassword(u.passwordHash)
	if err != nil {
		return nil, err
	}
	user, err := users.NewUser(
		id,
		name,
		email,
		passwordHash,
		u.isActive,
		u.createdAt,
		u.updatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (u *UsersStorage) GetByEmail(ctx context.Context, email users.Email) (*users.User, error) {
	log := u.log
	log.Info("attempting to get user by email", slog.String("email", email.String()))
	query := `SELECT id, name, email, password_hash, is_active, created_at, updated_at FROM users where email = $1;`
	var usr User
	err := u.db.Pool.QueryRow(ctx, query, email).Scan(&usr.id, &usr.name, &usr.email, &usr.passwordHash, &usr.isActive, &usr.createdAt, &usr.updatedAt)
	if err != nil {
		log.Error("failed to get user by email", slog.String("email", email.String()))
		return nil, err
	}

	log.Info("successful getting user by email", slog.String("email", email.String()))
	return mapperToDomain(usr)
}
