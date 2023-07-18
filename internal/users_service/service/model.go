package service

import (
	"github.com/google/uuid"
	"time"
)

type CreateUser struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUser struct {
	Username  *string   `json:"username"`
	Email     *string   `json:"email"`
	Password  *string   `json:"password"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SignInInput struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string
}
