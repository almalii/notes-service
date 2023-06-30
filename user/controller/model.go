package controller

import (
	"bookmarks/user/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,alpha,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (cur CreateUserRequest) ToDomain() (usecase.CreateUserInput, error) {
	validate := validator.New()
	err := validate.Struct(cur)
	if err != nil {
		return usecase.CreateUserInput{}, err.(validator.ValidationErrors)
	}

	return usecase.CreateUserInput{
		Username: cur.Username,
		Email:    cur.Email,
		Password: cur.Password,
	}, nil
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required,alpha,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (uur UpdateUserRequest) ToDomain() (usecase.UpdateUserInput, error) {
	validate := validator.New()
	err := validate.Struct(uur)
	if err != nil {
		return usecase.UpdateUserInput{}, err.(validator.ValidationErrors)
	}

	return usecase.UpdateUserInput{
		Username: &uur.Username,
		Email:    &uur.Email,
		Password: &uur.Password,
	}, nil
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (lr AuthRequest) ToDomain() (usecase.AuthInput, error) {
	validate := validator.New()
	err := validate.Struct(lr)
	if err != nil {
		return usecase.AuthInput{}, err.(validator.ValidationErrors)
	}

	return usecase.AuthInput{
		Email:    lr.Email,
		Password: lr.Password,
	}, nil
}

type UserResponseId struct {
	ID uuid.UUID `json:"id"`
}
