package v1

import (
	"context"
	pb_users_model "github.com/almalii/grpc-contracts/gen/go/users_service/model/v1"
	pb_users_service "github.com/almalii/grpc-contracts/gen/go/users_service/service/v1"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	usecase   UserUsecase
	validator *validator.Validate
	pb_users_service.UnimplementedUsersServiceServer
}

func (u *UsersServer) GetUser(
	ctx context.Context,
	req *pb_users_model.UserIDRequest,
) (*pb_users_model.GetUserResponse, error) {

	currentUserID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		logrus.Error("error getting user id from context")
		return nil, status.Error(codes.Internal, "error getting user id")
	}

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

	currentUserID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		logrus.Error("error getting user id from context")
		return nil, status.Error(codes.Internal, "error getting user id")
	}

	_, err := u.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error getting user: ", err)
		return nil, status.Error(codes.Internal, "error getting user")
	}

	input := dto.NewUpdateUserInput(req)

	if err = u.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = u.usecase.UpdateUser(ctx, input)
	if err != nil {
		logrus.Error("error updating user: ", err)
		return nil, status.Error(codes.Internal, "error updating user")
	}

	resp := dto.NewUpdateUserResponse(input)

	return resp, nil
}

func (u *UsersServer) DeleteUser(
	ctx context.Context,
	req *pb_users_model.UserIDRequest,
) (*emptypb.Empty, error) {

	currentUserID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		logrus.Error("error getting user id from context")
		return nil, status.Error(codes.Internal, "error getting user id")
	}

	_, err := u.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("id is not found: ", err)
		return nil, status.Error(codes.NotFound, "id is not found")
	}

	err = u.usecase.DeleteUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error deleting user: ", err)
		return nil, status.Error(codes.Internal, "error deleting user")
	}

	return &emptypb.Empty{}, nil
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
