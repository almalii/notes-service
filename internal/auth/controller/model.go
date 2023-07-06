package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/auth/usecase"
	"strings"
)

type SignUpResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func NewSignUpResponse(id uuid.UUID, username string, email string) SignUpResponse {
	return SignUpResponse{
		ID:       id,
		Username: username,
		Email:    email,
	}
}

type SignUpRequest struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=20"`
	Email    string `json:"email" validate:"required,email,min=5,max=254"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

func (sur SignUpRequest) ToDomain() (usecase.UserInput, error) {
	// TODO перенести создание валидатора в апп
	// TODO пернести валидацию в хендлер
	validate := validator.New()
	err := validate.Struct(sur)
	if err != nil {
		logrus.Error(err)
		return usecase.UserInput{}, err.(validator.ValidationErrors)
	}

	return usecase.UserInput{
		Username: sur.Username,
		Email:    strings.ToLower(sur.Email),
		Password: sur.Password,
	}, nil
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email,min=5,max=254"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

func (sir SignInRequest) ToDomain() (usecase.AuthInput, error) {
	validate := validator.New()
	err := validate.Struct(sir)
	if err != nil {
		logrus.Error(err)
		return usecase.AuthInput{}, err.(validator.ValidationErrors)
	}

	return usecase.AuthInput{
		Email:    strings.ToLower(sir.Email),
		Password: sir.Password,
	}, nil
}
