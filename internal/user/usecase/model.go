package usecase

import (
	"notes-rew/internal/user/service"
	"time"

	"github.com/google/uuid"
)

type UpdateUserInput struct {
	InitiatorID uuid.UUID
	Username    *string
	Email       *string
	Password    *string
}

func NewUpdateUserToService(username, email, password *string) service.UpdateUser {
	return service.UpdateUser{
		Username:  username,
		Email:     email,
		Password:  password,
		UpdatedAt: time.Now().UTC(),
	}
}

type AuthInput struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}
