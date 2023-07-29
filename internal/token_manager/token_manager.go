package token_manager

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

const (
	tokenTTL = 12 * time.Hour
)

type TokenManager interface {
	NewJWT(userID string) (string, error)
	ParseToken(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type tokenManager struct {
	signinKey string
}

func (t tokenManager) NewJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		Subject:   userID,
	})

	return token.SignedString([]byte(t.signinKey))
}

func (t tokenManager) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(t.signinKey), nil
	})
	if err != nil {
		return " uuid.Nil", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "uuid.Nil", errors.New("invalid token claims")
	}

	return claims.Subject, nil
}

func (t tokenManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func NewTokenManager(signinKey string) TokenManager {
	return &tokenManager{signinKey: signinKey}
}
