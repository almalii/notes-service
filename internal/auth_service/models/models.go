package models

import "github.com/google/uuid"

type AuthOutput struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
}

type AuthResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func NewAuthResponse(id uuid.UUID, username string, email string) AuthResponse {
	return AuthResponse{
		ID:       id,
		Username: username,
		Email:    email,
	}
}
