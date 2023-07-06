package usecase

import (
	"context"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/hash"
	"notes-rew/internal/user/models"
	"notes-rew/internal/user/service"

	"github.com/google/uuid"
)

type UserService interface {
	SaveUserByID(ctx context.Context, user service.CreateUser) error
	GetUserByID(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUserByID(ctx context.Context, id uuid.UUID, user service.UpdateUser) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	CheckerByEmail(ctx context.Context, email string) (bool, error)
}

type UserUsecase struct {
	service UserService
	hasher  hash.Hasher
}

func (u *UserUsecase) ReadUser(ctx context.Context, id uuid.UUID) (models.UserOutput, error) {
	return u.service.GetUserByID(ctx, id)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, req UpdateUserInput) error {
	hashedPassword, err := u.hasher.HasherPassword(*req.Password)
	if err != nil {
		logrus.Errorf("hash password error: %s", err)
	}

	// TODO check user can update user ( то есть это он )

	userUpdate := NewUpdateUserToService(req.Username, req.Email, &hashedPassword)

	err = u.service.UpdateUserByID(ctx, req.InitiatorID, userUpdate)
	if err != nil {
		logrus.Errorf("update user error: %s", err)
		return err
	}

	return nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.service.DeleteUserByID(ctx, id)
}

func (u *UserUsecase) CheckUserByEmail(ctx context.Context, email string) (bool, error) {
	return u.service.CheckerByEmail(ctx, email)
}

func NewUserUsecase(service UserService, hasher hash.Hasher) *UserUsecase {
	return &UserUsecase{
		service: service,
		hasher:  hasher,
	}
}
