package hash

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func HasherPassword(password string) (string, string, error) {
	salt := generateSalt()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatal("password hash error", err)
	}
	return string(hashedPassword), salt, nil
}

func ComparePassword(hashedPassword string, salt string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
}

func generateSalt() string {
	salt := make([]byte, 21)
	_, err := rand.Read(salt)
	if err != nil {
		logrus.Fatal("generate salt error", err)
	}
	return base64.URLEncoding.EncodeToString(salt)
}
