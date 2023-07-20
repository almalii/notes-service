package models

import "github.com/google/uuid"

type AuthOutput struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
}

type AuthResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func NewAuthResponse(id uuid.UUID, username string, email string) AuthResponse {
	return AuthResponse{
		UserID:   id,
		Username: username,
		Email:    email,
	}
}
