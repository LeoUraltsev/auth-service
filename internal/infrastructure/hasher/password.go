package hasher

import "golang.org/x/crypto/bcrypt"

type Password struct {
}

func NewHasher() *Password {
	return &Password{}
}

func (p Password) Verify(passwordHash []byte, password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(passwordHash, password)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p Password) Hash(password []byte) ([]byte, error) {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
