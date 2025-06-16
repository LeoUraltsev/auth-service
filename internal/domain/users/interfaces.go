package users

import "github.com/google/uuid"

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
	Verify(passwordHash []byte, password []byte) (bool, error)
}

type UserRepository interface {
	Save(user *User) error
	Get(id uuid.UUID) (*User, error)
	GetAll() ([]*User, error)
	ExistsByEmail(email Email) (bool, error)
}
