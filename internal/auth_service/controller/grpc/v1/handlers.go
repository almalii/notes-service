package v1

import (
	"context"
	pb_model "github.com/almalii/grpc-contracts/gen/go/auth_service/model/v1"
	pb_service "github.com/almalii/grpc-contracts/gen/go/auth_service/service/v1"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"notes-rew/internal/auth_service/models"
	"notes-rew/internal/auth_service/usecase"
)

type AuthUsecase interface {
	CreateUser(ctx context.Context, req usecase.UserInput) (uuid.UUID, error)
	AuthenticateUser(ctx context.Context, req usecase.AuthInput) (*models.AuthResponse, error)
}

type AuthServer struct {
	usecase   AuthUsecase
	validator *validator.Validate
	pb_service.UnimplementedAuthServiceServer
}

func (s *AuthServer) SignUp(ctx context.Context, req *pb_model.SignUpRequest) (*pb_model.SignUpResponse, error) {
	input := NewSignUpInput(req)

	if err := s.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID, err := s.usecase.CreateUser(ctx, input)
	if err != nil {
		logrus.Errorf("error creating user: %v", err)
		return nil, status.Error(codes.Internal, "error creating user")
	}

	resp := NewSignUpResponse(userID)

	return resp, nil
}

func (s *AuthServer) SignIn(ctx context.Context, req *pb_model.SignInRequest) (*pb_model.SignInResponse, error) {
	input := NewSignInInput(req)

	if err := s.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument")
	}

	authData, err := s.usecase.AuthenticateUser(ctx, input)
	if err != nil {
		logrus.Error("password is not correct")
		return nil, status.Errorf(codes.Unauthenticated, "password is not correct")
	}
	
	resp := NewSignInResponse(authData.Token)

	return resp, nil
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
