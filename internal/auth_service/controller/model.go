package controller

import (
	"github.com/google/uuid"
	"notes-rew/internal/auth_service/usecase"
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
	Email    string `json:"email" validate:"required,emailRFC,min=5,max=254"`
	Password string `json:"password" validate:"required,security"`
}

func (sur SignUpRequest) ToDomain() usecase.UserInput {
	return usecase.UserInput{
		Username: sur.Username,
		Email:    strings.ToLower(sur.Email),
		Password: sur.Password,
	}
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email,min=5,max=254"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

func (sir SignInRequest) ToDomain() usecase.AuthInput {
	return usecase.AuthInput{
		Email:    strings.ToLower(sir.Email),
		Password: sir.Password,
	}
}
