package handler

import (
	"context"
	"fmt"
	"notes-rew/internal/auth_service/controller/grpc/api"
)

type Server struct {
	api.UnimplementedAuthServiceServer
}

func NewServer(unimplementedAuthServiceServer api.UnimplementedAuthServiceServer) *Server {
	return &Server{UnimplementedAuthServiceServer: unimplementedAuthServiceServer}
}

func (s *Server) SignUp(context.Context, *api.SignUpRequest) (*api.SignUpResponse, error) {
	fmt.Println("SignUp")
	return &api.SignUpResponse{
		Id:       "1",
		Username: "2",
		Email:    "3",
	}, nil
}

func (s *Server) SignIn(context.Context, *api.SignInRequest) (*api.SignInResponse, error) {
	fmt.Println("SignIn")
	return &api.SignInResponse{
		Id:       "a",
		Username: "b",
		Email:    "c",
	}, nil
}
