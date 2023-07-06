package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"notes-rew/internal/user/usecase"
	"strings"
)

type UpdateUserRequest struct {
	Username *string `json:"username" validate:"required,alpha,min=3,max=20"`
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required"`
}

func (uur UpdateUserRequest) ToDomain(id uuid.UUID) (usecase.UpdateUserInput, error) {
	// TODO перенести создание валидатора в апп
	validate := validator.New()
	// TODO пернести валидацию в хендлер
	err := validate.Struct(uur)
	if err != nil {
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
