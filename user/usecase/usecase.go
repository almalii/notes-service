package usecase

import (
	"bookmarks/internal/hash"
	"bookmarks/user/models"
	"bookmarks/user/service"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	SaveUserByID(ctx context.Context, user service.CreateUser) error
	GetUserByID(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUserByID(ctx context.Context, id uuid.UUID, user service.UpdateUser) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	AuthByEmail(ctx context.Context, req service.SignInInput) (models.AuthResponse, error)
	CheckerByEmail(ctx context.Context, email string) (bool, error)
}

type UserUsecase struct {
	service UserService
}

func (u *UserUsecase) CreateUser(ctx context.Context, req CreateUserInput) (uuid.UUID, error) {
	hashedPassword, saltKey, err := hash.HasherPassword(req.Password)
	if err != nil {
		logrus.Errorf("hash password error: %s", err)
	}

	newUser := NewCreateUserInput(req.Username, req.Email, hashedPassword, saltKey)

	err = u.service.SaveUserByID(ctx, service.CreateUser(newUser))
	if err != nil {
		logrus.Errorf("save user error: %s", err)
		return uuid.Nil, err
	}

	return newUser.ID, nil
}

func (u *UserUsecase) ReadUser(ctx context.Context, id uuid.UUID) (models.UserOutput, error) {
	return u.service.GetUserByID(ctx, id)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserInput) error {
	hashedPassword, saltKey, err := hash.HasherPassword(*req.Password)
	if err != nil {
		logrus.Errorf("hash password error: %s", err)
	}

	userUpdate := NewUpdateUserInput(req.Username, req.Email, &hashedPassword, &saltKey)

	err = u.service.UpdateUserByID(ctx, id, service.UpdateUser(userUpdate))
	if err != nil {
		logrus.Errorf("update user error: %s", err)
		return err
	}

	return nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.service.DeleteUserByID(ctx, id)
}

func (u *UserUsecase) AuthenticateUser(ctx context.Context, req AuthInput) (models.AuthResponse, error) {

	return u.service.AuthByEmail(ctx, service.SignInInput(req))
}

func (u *UserUsecase) CheckUserByEmail(ctx context.Context, email string) (bool, error) {
	return u.service.CheckerByEmail(ctx, email)
}

func NewUserUsecase(service UserService) *UserUsecase {
	return &UserUsecase{
		service: service,
	}
}
