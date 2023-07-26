package v1

import (
	"context"
	"fmt"
	pb_users_model "github.com/almalii/grpc-contracts/gen/go/users_service/model/v1"
	pb_users_service "github.com/almalii/grpc-contracts/gen/go/users_service/service/v1"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"notes-rew/internal/users_service/models"
	"notes-rew/internal/users_service/models/dto"
	"notes-rew/internal/users_service/usecase"
	"strings"
)

type UserUsecase interface {
	ReadUser(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUser(ctx context.Context, req usecase.UpdateUserInput) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
}

type UsersServer struct {
	usecase   UserUsecase
	validator *validator.Validate
	pb_users_service.UnimplementedUsersServiceServer
}

func NewUsersServer(
	usecase UserUsecase,
	validator *validator.Validate,
	unimplementedUsersServiceServer pb_users_service.UnimplementedUsersServiceServer,
) *UsersServer {
	return &UsersServer{
		usecase:                         usecase,
		validator:                       validator,
		UnimplementedUsersServiceServer: unimplementedUsersServiceServer,
	}
}

func (u *UsersServer) GetUser(ctx context.Context, req *pb_users_model.UserIDRequest) (*pb_users_model.UserResponse, error) {
	currentUserID := dto.NewGetUserInput(req)

	user, err := u.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error getting user: ", err)
		return nil, err
	}

	resp := dto.NewGetUserResponse(user)

	return resp, nil
}

func (u *UsersServer) UpdateUser(ctx context.Context, req *pb_users_model.UpdateUserRequest) (*pb_users_model.UpdateUserResponse, error) {
	input := dto.NewUpdateUserInput(req)

	if err := u.validator.Struct(input); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		return nil, err
	}

	currentUserID := input.InitiatorID

	_, err := u.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error getting user: ", err)
		return nil, err
	}

	existingUser, err := u.usecase.CheckUserByEmail(ctx, strings.ToLower(*input.Email))
	if err != nil {
		logrus.Error("error checking user by email: ", err)
		return nil, err
	}

	if existingUser {
		logrus.Error("user with this email already exists")
		return nil, fmt.Errorf("user with this email already exists")
	}

	err = u.usecase.UpdateUser(ctx, input)
	if err != nil {
		logrus.Error("error updating user: ", err)
		return nil, err
	}

	resp := dto.NewUpdateUserResponse(input)

	return resp, nil
}

func (u *UsersServer) DeleteUser(ctx context.Context, req *pb_users_model.UserIDRequest) (*emptypb.Empty, error) {
	currentUserID := dto.NewDeleteUserInput(req)

	_, err := u.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("id is not found: ", err)
		return nil, err
	}

	err = u.usecase.DeleteUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error deleting user: ", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
