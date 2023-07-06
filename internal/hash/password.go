package hash

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/sirupsen/logrus"
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
		logrus.Fatal("password hash error", err)
	}
	return string(hashedPassword), nil
}

func (s *PasswordHasher) ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+s.salt))
}

func GenerateSalt() string {
	salt := make([]byte, 21)
	_, err := rand.Read(salt)
	if err != nil {
		logrus.Fatal("generate salt error", err)
	}
	return base64.URLEncoding.EncodeToString(salt)
}
