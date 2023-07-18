package service

import (
	"context"
	"github.com/google/uuid"
	"notes-rew/internal/users_service/models"
)

type UserStorage interface {
	CreateUserByID(ctx context.Context, user CreateUser) error
	GetUserByID(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUserByID(ctx context.Context, id uuid.UUID, user UpdateUser) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
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

func (s *UserService) CheckerByEmail(ctx context.Context, email string) (bool, error) {
	return s.storage.CheckUserByEmail(ctx, email)
}

func NewUserService(storage UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}
