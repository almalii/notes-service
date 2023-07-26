package v1

import (
	"context"
	"fmt"
	pb_users_model "github.com/almalii/grpc-contracts/gen/go/users_service/model/v1"
	pb_users_service "github.com/almalii/grpc-contracts/gen/go/users_service/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UsersServer struct {
	pb_users_service.UnimplementedUsersServiceServer
}

func NewUsersServer(unimplementedUsersServiceServer pb_users_service.UnimplementedUsersServiceServer) *UsersServer {
	return &UsersServer{UnimplementedUsersServiceServer: unimplementedUsersServiceServer}
}

func (u *UsersServer) GetUser(context.Context, *pb_users_model.UserIDRequest) (*pb_users_model.UserResponse, error) {
	fmt.Println("GetUser")
	return &pb_users_model.UserResponse{
		Id: "123",
	}, nil
}

func (u *UsersServer) UpdateUser(context.Context, *pb_users_model.UpdateUserRequest) (*pb_users_model.UpdateUserResponse, error) {
	fmt.Println("UpdateUser")
	return &pb_users_model.UpdateUserResponse{
		Id: "123",
	}, nil
}

func (u *UsersServer) DeleteUser(context.Context, *pb_users_model.UserIDRequest) (*emptypb.Empty, error) {
	fmt.Println("DeleteUser")
	return &emptypb.Empty{}, nil
}
