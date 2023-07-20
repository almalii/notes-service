package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"time"
)

const duration = 2 * time.Hour

type Session struct {
	ID        string
	Values    map[string]interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
	Options   *Options
}

type Options struct {
	HttpOnly bool
	Secure   bool
	SameSite string
}

// NewSession создает новый сеанс.
func NewSession() *Session {
	return &Session{
		ID:        GenerateUniqueID(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(duration),
		Values:    make(map[string]interface{}),
	}
}

func GenerateUniqueID() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		logrus.Error("error while generating random bytes")
		return err.Error()
	}
	// Преобразование случайных байтов в строку в формате UUID
	uuIdentification := hex.EncodeToString(bytes)

	return uuIdentification
}
