package users

import (
	"context"
	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password []byte) ([]byte, error)
}

type PasswordVerifier interface {
	Verify(passwordHash []byte, password []byte) (bool, error)
}

type UserServiceHandler interface {
	CreateUser(ctx context.Context, name string, email string, password string) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	GetListUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, name string, email string, password string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	Login(ctx context.Context, email string, password string) (string, error)
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email Email) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
}

type TokenGenerator interface {
	GenerateToken(userID uuid.UUID) (string, error)
}
