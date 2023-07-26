package dto

import (
	pb_model "github.com/almalii/grpc-contracts/gen/go/auth_service/model/v1"
	"github.com/google/uuid"
	"notes-rew/internal/auth_service/usecase"
	"strings"
)

func NewSignUpInput(req *pb_model.SignUpRequest) usecase.UserInput {
	return usecase.UserInput{
		Username: req.Username,
		Email:    strings.ToLower(req.Email),
		Password: req.Password,
	}
}

func NewSignUpResponse(id uuid.UUID) *pb_model.SignUpResponse {
	return &pb_model.SignUpResponse{
		Id: id.String(),
	}
}

func NewSignInInput(req *pb_model.SignInRequest) usecase.AuthInput {
	return usecase.AuthInput{
		Email:    strings.ToLower(req.Email),
		Password: req.Password,
	}
}

func NewSignInResponse(id uuid.UUID) *pb_model.SignInResponse {
	return &pb_model.SignInResponse{
		Id: id.String(),
	}
}
