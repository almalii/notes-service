package models

import (
	"github.com/google/uuid"
	"time"
)

type UserOutput struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthOutput struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}
