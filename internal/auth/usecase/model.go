package usecase

import (
	"github.com/google/uuid"
	"notes-rew/internal/auth/service"
	"time"
)

type UserInput struct {
	Username string
	Email    string
	Password string
}

type UserOutput struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuthInput struct {
	Email    string
	Password string
}

func NewUserOutput(username, email, passwordHash string) service.CreateUser {
	return service.CreateUser{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
}

func NewAuthInput(email, password string) service.SignInInput {
	return service.SignInInput{
		Email:    email,
		Password: password,
	}
}
