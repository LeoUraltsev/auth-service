package application

import (
	"context"
	"github.com/LeoUraltsev/auth-service/internal/domain/users"
	"github.com/LeoUraltsev/auth-service/internal/helper/logger"
	"github.com/google/uuid"
	"log/slog"
	"strings"
)

/*
todo: 1.для проверки существования и удаления нужно продумать транзакции так как используется несколько запросов к бд
есть риск невалидных данных. нужно продумать слой где будут запускаться транзакции

todo: 2. обработка ошибок, domain слоя

todo: 3. нужна обработка от возможности изменения и удаления чужих данных (сейчас по любому токену можно изменить любого юзера)
*/

type UserService interface {
	CreateUser(ctx context.Context, name string, email string, password string) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (*users.User, error)
	GetListUsers(ctx context.Context) ([]*users.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, name string, email string, password string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	Login(ctx context.Context, email string, password string) (string, error)
}

type UserServiceHandler struct {
	userRepo         users.UserRepository
	passwordHasher   users.PasswordHasher
	passwordVerifier users.PasswordVerifier
	tokenGen         users.TokenGenerator
	log              *slog.Logger
}

func NewUserService(
	userRepo users.UserRepository,
	passwordHasher users.PasswordHasher,
	passwordVerifier users.PasswordVerifier,
	tokenGen users.TokenGenerator,
	log *slog.Logger,
) *UserServiceHandler {
	return &UserServiceHandler{
		userRepo:         userRepo,
		passwordHasher:   passwordHasher,
		log:              log,
		passwordVerifier: passwordVerifier,
		tokenGen:         tokenGen,
	}
}

func (s *UserServiceHandler) CreateUser(ctx context.Context, name string, email string, password string) (uuid.UUID, error) {
	log := logger.LogWithContext(ctx, s.log)
	log.Info("creating user")

	n, err := users.NewName(name)
	if err != nil {
		log.Warn("failed to create user", slog.String("name", name), slog.String("error", err.Error()))
		return uuid.Nil, err
	}
	e, err := users.NewEmail(email)
	if err != nil {
		log.Warn("failed to create user", slog.String("email", email), slog.String("error", err.Error()))
		return uuid.Nil, err
	}
	if err = s.checkUniqueEmail(ctx, e); err != nil {
		log.Warn("failed to create user", slog.String("email", email), slog.String("error", err.Error()))
		return uuid.Nil, err
	}

	if strings.TrimSpace(password) == "" {
		log.Warn("failed to create user", slog.String("error", users.ErrPasswordRequired.Error()))
	}
	hashPassword, err := s.hashPassword([]byte(password))
	if err != nil {
		log.Warn("failed to create user", slog.String("error", err.Error()))
		return uuid.Nil, err
	}
	p, err := users.NewPassword(hashPassword)
	if err != nil {
		log.Warn("failed to create user", slog.String("error", err.Error()))
		return uuid.Nil, err
	}
	user, err := users.CreateUser(n, e, p)
	if err != nil {
		log.Warn("failed to create user", slog.Any("user", &user), slog.String("error", err.Error()))
		return uuid.Nil, err
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		log.Warn("failed to create user", slog.Any("user", user), slog.String("error", err.Error()))
		return uuid.Nil, err
	}
	log.Info("user created")
	return user.ID(), nil
}

func (s *UserServiceHandler) GetUser(ctx context.Context, id uuid.UUID) (*users.User, error) {
	log := logger.LogWithContext(ctx, s.log)
	log.Info("getting user")
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		log.Warn("failed to get user", slog.String("id", id.String()))
		return nil, err
	}
	log.Info("success getting user")
	return user, nil
}

func (s *UserServiceHandler) GetListUsers(ctx context.Context) ([]*users.User, error) {
	log := logger.LogWithContext(ctx, s.log)

	log.Info("getting all users")
	u, err := s.userRepo.GetAll(ctx)
	if err != nil {
		log.Warn("failed to get all users", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("success getting all users", slog.Any("users", &u))
	return u, nil
}

func (s *UserServiceHandler) UpdateUser(ctx context.Context, id uuid.UUID, name string, email string, password string) error {
	log := logger.LogWithContext(ctx, s.log)
	log.Info("updating user")

	log.Info("getting user")
	u, err := s.userRepo.Get(ctx, id)
	if err != nil {
		log.Warn("failed to get user", slog.String("id", id.String()))
		return err
	}

	if name != "" {
		newName, err := users.NewName(name)
		if err != nil {
			log.Warn("failed updating user name", slog.String("name", name), slog.String("error", err.Error()))
			return err
		}

		if err := u.UpdateName(newName); err != nil {
			log.Warn("failed to update user name", slog.String("name", name), slog.String("error", err.Error()))
			return err
		}
		log.Debug("success updating user name")
	}

	if email != "" {
		newEmail, err := users.NewEmail(email)
		if err != nil {
			log.Warn("failed to update user email", slog.String("email", email), slog.String("error", err.Error()))
			return err
		}
		if err := u.UpdateEmail(newEmail); err != nil {
			log.Warn("failed to update user email", slog.String("email", email), slog.String("error", err.Error()))
		}
		log.Debug("success updating user email")
	}

	if password != "" {
		hashPassword, err := s.hashPassword([]byte(password))
		if err != nil {
			log.Warn("failed to update user password", slog.String("error", err.Error()))
			return err
		}

		p, err := users.NewPassword(hashPassword)
		if err != nil {
			log.Warn("failed to update user password", slog.String("error", err.Error()))
			return err
		}

		if err := u.UpdatePassword(p); err != nil {
			log.Warn("failed to update user password", slog.String("error", err.Error()))
			return err
		}
		log.Debug("success updating user password")
	}

	err = s.userRepo.Save(ctx, u)
	if err != nil {
		log.Warn("failed to update user", slog.String("error", err.Error()))
		return err
	}
	log.Info("user updated")
	return nil
}

func (s *UserServiceHandler) DeleteUser(ctx context.Context, id uuid.UUID) error {
	log := logger.LogWithContext(ctx, s.log)
	log.Info("deleting user")
	u, err := s.userRepo.Get(ctx, id)
	if err != nil {
		log.Warn("failed to get user", slog.String("id", id.String()))
		return err
	}
	err = u.Delete()
	if err != nil {
		log.Warn("failed to delete user", slog.String("id", id.String()))
		return err
	}
	err = s.userRepo.Save(ctx, u)
	if err != nil {
		log.Warn("failed to delete user", slog.String("id", id.String()))
		return err
	}
	log.Info("user deleted")
	return nil
}

func (s *UserServiceHandler) Login(ctx context.Context, email string, password string) (string, error) {
	e, err := users.NewEmail(email)
	if err != nil {
		return "", err
	}
	p, err := users.NewPassword([]byte(password))
	if err != nil {
		return "", err
	}

	usr, err := s.userRepo.GetByEmail(ctx, e)
	if err != nil {
		return "", err
	}

	verify, err := s.passwordVerifier.Verify(usr.Password().Hash(), p.Hash())
	if err != nil {
		return "", err
	}
	if !verify {
		return "", users.ErrInvalidCredentials
	}
	token, err := s.tokenGen.GenerateToken(usr.ID())
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserServiceHandler) checkUniqueEmail(ctx context.Context, email users.Email) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return users.ErrEmailAlreadyExists
	}
	return nil
}

func (s *UserServiceHandler) hashPassword(password []byte) ([]byte, error) {
	return s.passwordHasher.Hash(password)
}
