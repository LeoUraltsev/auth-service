package application

import (
	"context"
	"errors"
	"github.com/LeoUraltsev/auth-service/internal/domain/users"
	"github.com/LeoUraltsev/auth-service/internal/helper/logger"
	"github.com/google/uuid"
	"log/slog"
	"strings"
)

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(repository users.UserRepository) error) error
}

type UserService interface {
	CreateUser(ctx context.Context, name string, email string, password string) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (*users.User, error)
	GetListUsers(ctx context.Context) ([]*users.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, name string, email string, password string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	Login(ctx context.Context, email string, password string) (string, error)
}

type UserServiceHandler struct {
	uof              UnitOfWork
	passwordHasher   users.PasswordHasher
	passwordVerifier users.PasswordVerifier
	tokenGen         users.TokenGenerator
	log              *slog.Logger
}

func NewUserService(
	uof UnitOfWork,
	passwordHasher users.PasswordHasher,
	passwordVerifier users.PasswordVerifier,
	tokenGen users.TokenGenerator,
	log *slog.Logger,
) *UserServiceHandler {
	return &UserServiceHandler{
		uof:              uof,
		passwordHasher:   passwordHasher,
		log:              log,
		passwordVerifier: passwordVerifier,
		tokenGen:         tokenGen,
	}
}

func (s *UserServiceHandler) CreateUser(ctx context.Context, name string, email string, password string) (uuid.UUID, error) {
	log := logger.LogWithContext(ctx, s.log)
	log.Info("creating user")
	var user *users.User
	err := s.uof.Execute(ctx, func(repo users.UserRepository) error {
		n, err := users.NewName(name)
		if err != nil {
			log.Warn("failed to create user", slog.String("name", name), slog.String("error", err.Error()))
			return err
		}
		e, err := users.NewEmail(email)
		if err != nil {
			log.Warn("failed to create user", slog.String("email", email), slog.String("error", err.Error()))
			return err
		}
		if err = s.checkUniqueEmail(ctx, e); err != nil {
			log.Warn("failed to create user", slog.String("email", email), slog.String("error", err.Error()))
			return err
		}

		if strings.TrimSpace(password) == "" {
			log.Warn("failed to create user", slog.String("error", users.ErrPasswordRequired.Error()))
		}
		hashPassword, err := s.hashPassword([]byte(password))
		if err != nil {
			log.Warn("failed to create user", slog.String("error", err.Error()))
			return err
		}
		p, err := users.NewPassword(hashPassword)
		if err != nil {
			log.Warn("failed to create user", slog.String("error", err.Error()))
			return err
		}
		user, err = users.CreateUser(n, e, p)
		if err != nil {
			log.Warn("failed to create user", slog.Any("user", &user), slog.String("error", err.Error()))
			return err
		}

		if err := repo.Save(ctx, user); err != nil {
			log.Warn("failed to create user", slog.Any("user", user), slog.String("error", err.Error()))
			return err
		}
		log.Info("user created")
		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID(), nil
}

func (s *UserServiceHandler) GetUser(ctx context.Context, id uuid.UUID) (*users.User, error) {
	log := logger.LogWithContext(ctx, s.log)
	log.Info("getting user")

	var user *users.User
	err := s.uof.Execute(ctx, func(repo users.UserRepository) error {
		var err error
		user, err = repo.Get(ctx, id)
		if err != nil {
			log.Warn("failed to get user", slog.String("id", id.String()))
			return err
		}
		log.Info("success getting user")
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceHandler) GetListUsers(ctx context.Context) ([]*users.User, error) {
	log := logger.LogWithContext(ctx, s.log)
	var u []*users.User
	log.Info("getting all users")

	err := s.uof.Execute(ctx, func(repo users.UserRepository) error {
		u, err := repo.GetAll(ctx)
		if err != nil {
			log.Warn("failed to get all users", slog.String("error", err.Error()))
			return err
		}
		log.Info("success getting all users", slog.Any("users", &u))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserServiceHandler) UpdateUser(ctx context.Context, id uuid.UUID, name string, email string, password string) error {
	log := logger.LogWithContext(ctx, s.log)
	log.Info("updating user")

	var u *users.User

	err := s.uof.Execute(ctx, func(repo users.UserRepository) error {
		ctxUserID, err := userIDFromContext(ctx)
		if err != nil {
			log.Error("user id missing in context", slog.String("id", id.String()))
			return err
		}

		if ctxUserID != id {
			log.Error("user id mismatch token user", slog.String("id", id.String()))
			return errors.New("user id mismatch token user")
		}

		log.Info("getting user")
		u, err = repo.Get(ctx, id)
		if err != nil {
			log.Warn("failed to get user", slog.String("id", id.String()))
			return err
		}

		if !u.IsActive() {
			log.Warn("user isnt active", slog.String("id", id.String()))
			return errors.New("user isnt active")
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

		err = repo.Save(ctx, u)
		if err != nil {
			log.Warn("failed to update user", slog.String("error", err.Error()))
			return err
		}
		log.Info("user updated")
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *UserServiceHandler) DeleteUser(ctx context.Context, id uuid.UUID) error {
	log := logger.LogWithContext(ctx, s.log)
	log.Info("deleting user")

	var u *users.User

	err := s.uof.Execute(ctx, func(repo users.UserRepository) error {
		ctxUserID, err := userIDFromContext(ctx)
		if err != nil {
			log.Error("user id missing in context", slog.String("id", id.String()))
			return err
		}

		if ctxUserID != id {
			log.Error("user id mismatch token user", slog.String("id", id.String()))
			return errors.New("user id mismatch token user")
		}

		u, err = repo.Get(ctx, id)
		if err != nil {
			log.Warn("failed to get user", slog.String("id", id.String()))
			return err
		}

		if !u.IsActive() {
			log.Warn("user isnt active", slog.String("id", id.String()))
			return errors.New("user isnt active")
		}

		err = u.Delete()
		if err != nil {
			log.Warn("failed to delete user", slog.String("id", id.String()))
			return err
		}
		err = repo.Save(ctx, u)
		if err != nil {
			log.Warn("failed to delete user", slog.String("id", id.String()))
			return err
		}
		log.Info("user deleted")
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *UserServiceHandler) Login(ctx context.Context, email string, password string) (string, error) {
	var token string
	err := s.uof.Execute(ctx, func(repo users.UserRepository) error {
		e, err := users.NewEmail(email)
		if err != nil {
			return err
		}
		p, err := users.NewPassword([]byte(password))
		if err != nil {
			return err
		}

		usr, err := repo.GetByEmail(ctx, e)
		if err != nil {
			return err
		}

		verify, err := s.passwordVerifier.Verify(usr.Password().Hash(), p.Hash())
		if err != nil {
			return err
		}
		if !verify {
			return users.ErrInvalidCredentials
		}
		token, err = s.tokenGen.GenerateToken(usr.ID())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserServiceHandler) checkUniqueEmail(ctx context.Context, email users.Email) error {
	err := s.uof.Execute(ctx, func(repo users.UserRepository) error {
		exists, err := repo.ExistsByEmail(ctx, email)
		if err != nil {
			return err
		}
		if exists {
			return users.ErrEmailAlreadyExists
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *UserServiceHandler) hashPassword(password []byte) ([]byte, error) {
	return s.passwordHasher.Hash(password)
}

// todo: возможно нужен пакет для хранения констант-ключей
func userIDFromContext(ctx context.Context) (uuid.UUID, error) {
	u, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user id not found in context")
	}
	return u, nil
}
