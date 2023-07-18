package service

import (
	"github.com/google/uuid"
	"time"
)

type CreateUser struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SignInInput struct {
	Email    string
	Password string
}
