package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"time"
)

const duration = 24 * time.Hour

type Session struct {
	ID        string
	Values    map[string]interface{}
	CreatedAt time.Time
	ExpiresAt time.Duration
}

func NewSession() *Session {
	return &Session{
		ID:        generateToken(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: duration,
		Values:    make(map[string]interface{}),
	}
}

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		logrus.Error("error while generating random bytes")
		return err.Error()
	}

	return base64.RawURLEncoding.EncodeToString(b)
}
