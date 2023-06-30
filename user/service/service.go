package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/hash"
	"notes-rew/user/models"
)

type UserStorage interface {
	CreateUserByID(ctx context.Context, user CreateUser) error
	GetUserByID(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUserByID(ctx context.Context, id uuid.UUID, user UpdateUser) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	GetUserForAuth(ctx context.Context, email string) (models.AuthOutput, error)
	CheckerByEmail(ctx context.Context, email string) (bool, error)
}

type UserService struct {
	storage UserStorage
}

func (s *UserService) SaveUserByID(ctx context.Context, user CreateUser) error {
	return s.storage.CreateUserByID(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (models.UserOutput, error) {
	return s.storage.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUserByID(ctx context.Context, id uuid.UUID, user UpdateUser) error {
	return s.storage.UpdateUserByID(ctx, id, user)
}

func (s *UserService) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	return s.storage.DeleteUserByID(ctx, id)
}

func (s *UserService) AuthByEmail(ctx context.Context, req SignInInput) (models.AuthResponse, error) {
	user, err := s.storage.GetUserForAuth(ctx, req.Email)

	if err != nil {
		logrus.Errorf("user not found: %s", err)
		return models.AuthResponse{}, fmt.Errorf("user not found")
	}
	if user == (models.AuthOutput{}) {
		logrus.Errorf("user not found")
		return models.AuthResponse{}, fmt.Errorf("user not found")
	}
	if err = hash.ComparePassword(user.Password, user.SaltKey, req.Password); err != nil {
		logrus.Errorf("password is not correct")
		return models.AuthResponse{}, fmt.Errorf("password is not correct")
	}

	return models.NewAuthResponse(user), nil
}

func (s *UserService) CheckerByEmail(ctx context.Context, email string) (bool, error) {
	return s.storage.CheckerByEmail(ctx, email)
}

func NewUserService(storage UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}
