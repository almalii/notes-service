package controller

import (
	"github.com/google/uuid"
	"notes-rew/internal/users_service/usecase"
	"strings"
)

type UpdateUserRequest struct {
	Username *string `json:"username" validators:"required,alphanum,min=3,max=20"`
	Email    *string `json:"email" validate:"required,emailRFC,min=5,max=254"`
	Password *string `json:"password" validate:"required,security"`
}

func (uur UpdateUserRequest) ToDomain(id uuid.UUID) usecase.UpdateUserInput {
	emailToLower := strings.ToLower(*uur.Email)

	return usecase.UpdateUserInput{
		InitiatorID: id,
		Username:    uur.Username,
		Email:       &emailToLower,
		Password:    uur.Password,
	}
}
