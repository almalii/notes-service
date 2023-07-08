package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"notes-rew/internal/user/usecase"
	"strings"
)

type UpdateUserRequest struct {
	Username *string `json:"username" validators:"required,alphanum,min=3,max=20"`
	Email    *string `json:"email" validate:"required,emailRFC,min=5,max=254"`
	Password *string `json:"password" validate:"required,security"`
}

func (uur UpdateUserRequest) ToDomain(id uuid.UUID, validate *validator.Validate) (usecase.UpdateUserInput, error) {
	// TODO пернести валидацию в хендлер
	if err := validate.Struct(uur); err != nil {
		return usecase.UpdateUserInput{}, err.(validator.ValidationErrors)
	}
	emailToLower := strings.ToLower(*uur.Email)

	return usecase.UpdateUserInput{
		InitiatorID: id,
		Username:    uur.Username,
		Email:       &emailToLower,
		Password:    uur.Password,
	}, nil
}
