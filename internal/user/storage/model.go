package storage

import (
	"github.com/google/uuid"
	"notes-rew/internal/user/models"
	"time"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}

type SignInResponse struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

func NewAuthResponse(resp SignInResponse) models.AuthOutput {
	return models.AuthOutput{
		ID:       resp.ID,
		Password: resp.Password,
	}
}
