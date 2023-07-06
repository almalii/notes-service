package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/auth/models"
	"notes-rew/internal/auth/service"
	"notes-rew/internal/hash"
)

type AuthService interface {
	CreateUserServ(ctx context.Context, user service.CreateUser) error
	AuthByEmail(ctx context.Context, req service.SignInInput) (models.AuthOutput, error)
	CheckerByEmail(ctx context.Context, email string) (bool, error)
}

type AuthUsecase struct {
	service AuthService
}

func (u *AuthUsecase) CreateUser(ctx context.Context, req UserInput) (uuid.UUID, error) {
	hashedPassword, err := hash.HasherPassword(req.Password)
	if err != nil {
		logrus.Errorf("hash password error: %s", err)
	}

	newUser := NewUserOutput(req.Username, req.Email, hashedPassword)

	err = u.service.CreateUserServ(ctx, newUser)
	if err != nil {
		logrus.Errorf("save user error: %s", err)
		return uuid.Nil, err
	}

	return newUser.ID, nil
}

func (u *AuthUsecase) AuthenticateUser(ctx context.Context, req AuthInput) (models.AuthResponse, error) {
	user, err := u.service.AuthByEmail(ctx, service.SignInInput(req))
	if err != nil {
		logrus.Errorf("user not found: %s", err)
		return models.AuthResponse{}, err
	}

	if err = hash.ComparePassword(user.PasswordHash, req.Password); err != nil {
		logrus.Errorf("password is not correct: %s", err)
		return models.AuthResponse{}, err
	}
	return models.NewAuthResponse(user.ID, user.Username, user.Email), nil
}

func (u *AuthUsecase) CheckUserByEmail(ctx context.Context, email string) (bool, error) {
	return u.service.CheckerByEmail(ctx, email)
}

func NewAuthUsecase(service AuthService) *AuthUsecase {
	return &AuthUsecase{service: service}
}
