package usecase

import (
	"notes-rew/internal/users_service/service"
	"time"

	"github.com/google/uuid"
)

type UpdateUserInput struct {
	InitiatorID uuid.UUID
	Username    *string `json:"username" validators:"required,alphanum,min=3,max=20"`
	Email       *string `json:"email" validate:"required,emailRFC,min=5,max=254"`
	Password    *string `json:"password" validate:"required,security"`
}

func NewUpdateUserToService(username, email, password *string) service.UpdateUser {
	return service.UpdateUser{
		Username:  username,
		Email:     email,
		Password:  password,
		UpdatedAt: time.Now().UTC(),
	}
}
