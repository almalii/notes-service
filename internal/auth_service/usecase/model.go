package usecase

import (
	"github.com/google/uuid"
	"notes-rew/internal/auth_service/service"
	"time"
)

type UserInput struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=20"`
	Email    string `json:"email" validate:"required,emailRFC,min=5,max=254"`
	Password string `json:"password" validate:"required,security"`
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
	Email    string `json:"email" validate:"required,email,min=5,max=254"`
	Password string `json:"password" validate:"required,min=6,max=30"`
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
