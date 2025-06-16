package users

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrEmailNotValid      = errors.New("email is not valid")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNameRequired       = errors.New("name is required")
	ErrPasswordRequired   = errors.New("password is required")
)

type Name string
type Email struct {
	value string
}
type Password struct {
	hash []byte
}
type User struct {
	id        uuid.UUID
	name      Name
	email     Email
	password  Password
	isActive  bool
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(
	id uuid.UUID,
	name Name,
	email Email,
	passwordHash Password,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	if err := email.validate(); err != nil {
		return nil, err
	}
	if err := name.validate(); err != nil {
		return nil, err
	}
	if err := passwordHash.validate(); err != nil {
		return nil, err
	}
	return &User{
		id:        id,
		name:      name,
		email:     email,
		password:  passwordHash,
		isActive:  isActive,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func CreateUser(
	name Name,
	email Email,
	password Password,
) (*User, error) {
	id := uuid.New()
	return NewUser(id, name, email, password, true, time.Now().UTC(), time.Now().UTC())
}

func (u *User) ID() uuid.UUID {
	return u.id
}
func (u *User) Name() Name {
	return u.name
}
func (u *User) Email() Email {
	return u.email
}
func (u *User) Password() Password {
	return u.password
}
func (u *User) IsActive() bool {
	return u.isActive
}
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) UpdateEmail(email Email) error {
	err := email.validate()
	if err != nil {
		return err
	}
	u.email = email
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) UpdatePassword(password Password) error {
	err := password.validate()
	if err != nil {
		return err
	}
	u.password = password
	u.updatedAt = time.Now().UTC()
	return nil
}

func NewEmail(email string) Email {
	return Email{
		value: email,
	}
}

func (e Email) validate() error {
	if len(e.value) == 0 {
		return ErrEmailNotValid
	}
	return nil
}

func (e Email) String() string {
	return e.value
}

func NewPassword(passwordHash string) Password {
	return Password{
		hash: []byte(passwordHash),
	}
}

func (p Password) String() string {
	return string(p.hash)
}

func (p Password) validate() error {
	if len(p.hash) == 0 {
		return ErrPasswordRequired
	}
	return nil
}

func NewName(name string) Name {
	return Name(name)
}

func (n Name) validate() error {
	if n == "" {
		return ErrNameRequired
	}
	return nil
}

func (n Name) String() string {
	return n.String()
}
