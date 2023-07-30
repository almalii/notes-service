package service

import (
	"context"
	"notes-rew/internal/auth_service/models"
)

type AuthStorage interface {
	SaveUserToDB(ctx context.Context, user CreateUser) error
	GetUserForAuth(ctx context.Context, email string) (models.AuthOutput, error)
	CheckUserByEmail(ctx context.Context, email string) error
}

type AuthService struct {
	storage AuthStorage
}

func (s *AuthService) CreateUserServ(ctx context.Context, user CreateUser) error {
	return s.storage.SaveUserToDB(ctx, user)
}

func (s *AuthService) CheckUserByEmail(ctx context.Context, email string) error {
	return s.storage.CheckUserByEmail(ctx, email)
}

func (s *AuthService) AuthByEmail(ctx context.Context, req SignInInput) (models.AuthOutput, error) {
	return s.storage.GetUserForAuth(ctx, req.Email)
}

func NewAuthService(storage AuthStorage) *AuthService {
	return &AuthService{storage: storage}
}
