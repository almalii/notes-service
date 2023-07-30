package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/hash"
	"notes-rew/internal/users_service/models"
	"notes-rew/internal/users_service/service"
	"strings"

	"github.com/google/uuid"
)

type UserService interface {
	SaveUserByID(ctx context.Context, user service.CreateUser) error
	GetUserByID(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUserByID(ctx context.Context, id uuid.UUID, user service.UpdateUser) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	CheckerByEmail(ctx context.Context, email string) error
}

type UserUsecase struct {
	service   UserService
	hasher    hash.Hasher
	validator *validator.Validate
}

func (u *UserUsecase) ReadUser(ctx context.Context, id uuid.UUID) (models.UserOutput, error) {
	return u.service.GetUserByID(ctx, id)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, req UpdateUserInput) error {
	if err := u.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		return err
	}

	err := u.service.CheckerByEmail(ctx, strings.ToLower(*req.Email))
	if err != nil {
		return err
	}

	hashedPassword, err := u.hasher.HasherPassword(*req.Password)
	if err != nil {
		logrus.Errorf("hash password error: %s", err)
	}

	userUpdate := NewUpdateUserToService(req.Username, req.Email, &hashedPassword)

	err = u.service.UpdateUserByID(ctx, req.InitiatorID, userUpdate)
	if err != nil {
		logrus.Errorf("update users error: %s", err)
		return err
	}

	return nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.service.DeleteUserByID(ctx, id)
}

func NewUserUsecase(service UserService, hasher hash.Hasher, validator *validator.Validate) *UserUsecase {
	return &UserUsecase{
		service:   service,
		hasher:    hasher,
		validator: validator,
	}
}
