package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/auth_service/models"
	"notes-rew/internal/auth_service/service"
	"notes-rew/internal/hash"
	"notes-rew/internal/token_manager"
	"strings"
)

type AuthService interface {
	CreateUserServ(ctx context.Context, user service.CreateUser) error
	AuthByEmail(ctx context.Context, req service.SignInInput) (models.AuthOutput, error)
	CheckUserByEmail(ctx context.Context, email string) error
}

type AuthUsecase struct {
	service      AuthService
	hasher       hash.Hasher
	tokenManager *token_manager.TokenManager
}

func (u *AuthUsecase) CreateUser(ctx context.Context, req UserInput) (uuid.UUID, error) {
	err := u.service.CheckUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		logrus.Printf("check user error: %s", err)
		return uuid.Nil, err
	}

	hashedPassword, err := u.hasher.HasherPassword(req.Password)
	if err != nil {
		logrus.Printf("hash password error: %s", err)
		return uuid.Nil, err
	}

	newUser := NewUserOutput(req.Username, req.Email, hashedPassword) // возвращает новый id

	err = u.service.CreateUserServ(ctx, newUser)
	if err != nil {
		logrus.Printf("save users error: %s", err)
		return uuid.Nil, err
	}

	return newUser.ID, nil
}

func (u *AuthUsecase) AuthenticateUser(ctx context.Context, req AuthInput) (*models.AuthResponse, error) {
	user, err := u.service.AuthByEmail(ctx, service.SignInInput(req))
	if err != nil {
		logrus.Printf("user not found: %s", err)
		return nil, err
	}

	if err = u.hasher.ComparePassword(user.PasswordHash, req.Password); err != nil {
		logrus.Printf("password is not correct: %s", err)
		return nil, err
	}

	jwt, err := u.tokenManager.NewJWT(user.UserID.String())
	if err != nil {
		logrus.Printf("jwt error: %s", err)
		return nil, err
	}

	resp := NewAuthResponse(jwt)

	return resp, nil
}

func NewAuthUsecase(service AuthService, hasher hash.Hasher, tokenManager *token_manager.TokenManager) *AuthUsecase {
	return &AuthUsecase{
		service:      service,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}
