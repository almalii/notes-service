package v1

import (
	"context"
	"fmt"
	pb_model "github.com/almalii/grpc-contracts/gen/go/auth_service/model/v1"
	pb_service "github.com/almalii/grpc-contracts/gen/go/auth_service/service/v1"
)

type AuthServer struct {
	pb_service.UnimplementedAuthServiceServer
}

func NewAuthServer(unimplementedAuthServiceServer pb_service.UnimplementedAuthServiceServer) *AuthServer {
	return &AuthServer{UnimplementedAuthServiceServer: unimplementedAuthServiceServer}
}

func (s *AuthServer) SignUp(context.Context, *pb_model.SignUpRequest) (*pb_model.SignUpResponse, error) {
	fmt.Println("SignUp")
	return &pb_model.SignUpResponse{
		Id: "123",
	}, nil
}

func (s *AuthServer) SignIn(context.Context, *pb_model.SignInRequest) (*pb_model.SignInResponse, error) {
	fmt.Println("SignIn")
	return &pb_model.SignInResponse{
		Id: "123",
	}, nil
}
