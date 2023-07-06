package hash

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"notes-rew/internal/config"
)

func HasherPassword(password string) (string, error) {
	salt := config.SaltKey()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatal("password hash error", err)
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword string, password string) error {
	salt := config.SaltKey()
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
}

//func generateSalt() string {
//	salt := make([]byte, 21)
//	_, err := rand.Read(salt)
//	if err != nil {
//		logrus.Fatal("generate salt error", err)
//	}
//	return base64.URLEncoding.EncodeToString(salt)
//}
