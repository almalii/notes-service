package v1

import (
	"context"
	"errors"
	pb_model "github.com/almalii/grpc-contracts/gen/go/auth_service/model/v1"
	pb_service "github.com/almalii/grpc-contracts/gen/go/auth_service/service/v1"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/auth_service/models"
	"notes-rew/internal/auth_service/models/dto"
	"notes-rew/internal/auth_service/usecase"
	"strings"
)

type AuthUsecase interface {
	CreateUser(ctx context.Context, req usecase.UserInput) (uuid.UUID, error)
	AuthenticateUser(ctx context.Context, req usecase.AuthInput) (models.AuthResponse, error)
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
}

type AuthServer struct {
	usecase   AuthUsecase
	validator *validator.Validate
	pb_service.UnimplementedAuthServiceServer
}

func NewAuthServer(
	usecase AuthUsecase,
	validator *validator.Validate,
	unimplementedAuthServiceServer pb_service.UnimplementedAuthServiceServer,
) *AuthServer {
	return &AuthServer{
		usecase:                        usecase,
		validator:                      validator,
		UnimplementedAuthServiceServer: unimplementedAuthServiceServer,
	}
}

func (s *AuthServer) SignUp(ctx context.Context, req *pb_model.SignUpRequest) (*pb_model.SignUpResponse, error) {
	input := dto.NewSignUpInput(req)

	existingUser, err := s.usecase.CheckUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		logrus.Errorf("no such user exists")
		return nil, errors.New("no such user exists")
	}

	if existingUser {
		logrus.Error("email already exists")
		return nil, errors.New("email already exists")
	}

	if err = s.validator.Struct(input); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		return nil, errors.New("validator error")
	}

	userID, err := s.usecase.CreateUser(ctx, input)
	if err != nil {
		logrus.Errorf("error creating user: %v", err)
		return nil, err
	}

	resp := dto.NewSignUpResponse(userID)

	return resp, nil
}

func (s *AuthServer) SignIn(ctx context.Context, req *pb_model.SignInRequest) (*pb_model.SignInResponse, error) {
	input := dto.NewSignInInput(req)

	if err := s.validator.Struct(input); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		return nil, errors.New("validator error")
	}

	resp, err := s.usecase.AuthenticateUser(ctx, input)
	if err != nil {
		logrus.Error("password is not correct")
		return nil, errors.New("password is not correct")
	}

	respID := dto.NewSignInResponse(resp.UserID)

	return respID, nil
}
