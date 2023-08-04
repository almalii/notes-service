package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HasherPassword(password string) (string, error)
	ComparePassword(hashedPassword string, password string) error
}

type PasswordHasher struct {
	salt string
}

func NewPasswordHasher(salt string) *PasswordHasher {
	return &PasswordHasher{
		salt: salt,
	}
}

func (s *PasswordHasher) HasherPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+s.salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *PasswordHasher) ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+s.salt))
}
