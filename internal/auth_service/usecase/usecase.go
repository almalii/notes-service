package usecase

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/auth_service/models"
	"notes-rew/internal/auth_service/service"
	"notes-rew/internal/hash"
	"notes-rew/internal/token_manager"
	"strings"
	"time"
)

type AuthService interface {
	CreateUserServ(ctx context.Context, user service.CreateUser) error
	AuthByEmail(ctx context.Context, req service.SignInInput) (models.AuthOutput, error)
}

const (
	tokenTTL = 12 * time.Hour
)

type AuthUsecase struct {
	service      AuthService
	hasher       hash.Hasher
	tokenManager token_manager.TokenManager
	validator    *validator.Validate
}

func (u *AuthUsecase) CreateUser(ctx context.Context, req UserInput) (uuid.UUID, error) {
	if err := u.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		return uuid.Nil, err
	}

	user, err := u.service.AuthByEmail(ctx, service.SignInInput{
		Email: strings.ToLower(req.Email),
	})

	if err != nil {
		logrus.Errorf("users_service not found: %s", err)
		return uuid.Nil, err
	}

	if user.Email == req.Email {
		return uuid.Nil, errors.New("email already exists")
	}

	hashedPassword, err := u.hasher.HasherPassword(req.Password)
	if err != nil {
		logrus.Errorf("hash password error: %s", err)
	}

	newUser := NewUserOutput(req.Username, req.Email, hashedPassword)

	err = u.service.CreateUserServ(ctx, newUser)
	if err != nil {
		logrus.Errorf("save users_service error: %s", err)
		return uuid.Nil, err
	}

	return newUser.ID, nil
}

func (u *AuthUsecase) AuthenticateUser(ctx context.Context, req AuthInput) (*models.AuthResponse, error) {
	if err := u.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		return nil, err
	}

	user, err := u.service.AuthByEmail(ctx, service.SignInInput(req))
	if err != nil {
		logrus.Errorf("users_service not found: %s", err)
		return nil, err
	}

	if err = u.hasher.ComparePassword(user.PasswordHash, req.Password); err != nil {
		logrus.Errorf("password is not correct: %s", err)
		return nil, err
	}

	jwt, err := u.tokenManager.NewJWT(user.UserID.String())
	if err != nil {
		logrus.Errorf("jwt error: %s", err)
		return nil, err
	}

	resp := models.AuthResponse{
		Token: jwt,
	}

	return &resp, nil
}

func NewAuthUsecase(service AuthService, hasher hash.Hasher, tokenManager token_manager.TokenManager, validator *validator.Validate) *AuthUsecase {
	return &AuthUsecase{
		service:      service,
		hasher:       hasher,
		tokenManager: tokenManager,
		validator:    validator,
	}
}
