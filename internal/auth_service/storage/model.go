package storage

import (
	"github.com/google/uuid"
	"notes-rew/internal/auth_service/models"
	"time"
)

type SaveUser struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuthResponse struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
}

func NewAuthResponse(id uuid.UUID, username string, email string, passwordHash string) models.AuthOutput {
	return models.AuthOutput{
		UserID:       id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}
}
