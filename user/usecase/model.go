package usecase

import (
	"github.com/google/uuid"
	"time"
)

type CreateUserInput struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	SaltKey   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCreateUserInput(username string, email string, password string, salt string) CreateUserInput {
	return CreateUserInput{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  password,
		SaltKey:   salt,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

type UpdateUserInput struct {
	Username  *string
	Email     *string
	Password  *string
	SaltKey   *string
	UpdatedAt time.Time
}

func NewUpdateUserInput(username *string, email *string, password *string, salt *string) UpdateUserInput {
	return UpdateUserInput{
		Username:  username,
		Email:     email,
		Password:  password,
		SaltKey:   salt,
		UpdatedAt: time.Now().UTC(),
	}
}

type AuthInput struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
	SaltKey  string
}
