package v1

import (
	"context"
	pb_users_model "github.com/almalii/grpc-contracts/gen/go/users_service/model/v1"
	pb_users_service "github.com/almalii/grpc-contracts/gen/go/users_service/service/v1"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"notes-rew/internal/users_service/models"
	"notes-rew/internal/users_service/models/dto"
	"notes-rew/internal/users_service/usecase"
)

type UserUsecase interface {
	ReadUser(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUser(ctx context.Context, req usecase.UpdateUserInput) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UsersServer struct {
	usecase UserUsecase
	pb_users_service.UnimplementedUsersServiceServer
}

func (u *UsersServer) GetUser(
	ctx context.Context,
	req *pb_users_model.UserIDRequest,
) (*pb_users_model.GetUserResponse, error) {
	currentUserID := ctx.Value("userID").(uuid.UUID)

	user, err := u.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error getting user: ", err)
		return nil, err
	}

	resp := dto.NewGetUserResponse(user)

	return resp, nil
}

func (u *UsersServer) UpdateUser(
	ctx context.Context,
	req *pb_users_model.UpdateUserRequest,
) (*pb_users_model.UpdateUserResponse, error) {
	currentUserID := ctx.Value("userID").(uuid.UUID)

	_, err := u.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error getting user: ", err)
		return nil, err
	}

	input := dto.NewUpdateUserInput(req)

	err = u.usecase.UpdateUser(ctx, input)
	if err != nil {
		logrus.Error("error updating user: ", err)
		return nil, err
	}

	resp := dto.NewUpdateUserResponse(input)

	return resp, nil
}

func (u *UsersServer) DeleteUser(
	ctx context.Context,
	req *pb_users_model.UserIDRequest,
) (*emptypb.Empty, error) {
	currentUserID := ctx.Value("userID").(uuid.UUID)

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

func NewUsersServer(
	usecase UserUsecase,
	unimplementedUsersServiceServer pb_users_service.UnimplementedUsersServiceServer,
) *UsersServer {
	return &UsersServer{
		usecase:                         usecase,
		UnimplementedUsersServiceServer: unimplementedUsersServiceServer,
	}
}
