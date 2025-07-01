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
