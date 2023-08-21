package controller

import (
	"context"

	"github.com/almalii/swagger-contracts/restapi/operations"
	"github.com/almalii/swagger-contracts/restapi/operations/auth"
	middleware2 "github.com/go-openapi/runtime/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"notes-rew/internal/auth_service/models"
	"notes-rew/internal/auth_service/usecase"
)

type AuthUsecase interface {
	CreateUser(ctx context.Context, req usecase.UserInput) (uuid.UUID, error)
	AuthenticateUser(ctx context.Context, req usecase.AuthInput) (*models.AuthResponse, error)
}

type AuthController struct {
	usecase   AuthUsecase
	validator *validator.Validate
}

func (c *AuthController) Register(api *operations.NotesAPIAPI) {
	api.AuthPostAuthLoginHandler = auth.PostAuthLoginHandlerFunc(c.AuthLogin)
	api.AuthPostAuthRegisterHandler = auth.PostAuthRegisterHandlerFunc(c.AuthRegister)
}

func (c *AuthController) AuthLogin(params auth.PostAuthLoginParams) middleware2.Responder {
	//
	return auth.NewPostAuthLoginOK()
}

func (c *AuthController) AuthRegister(params auth.PostAuthRegisterParams) middleware2.Responder {
	//
	return auth.NewPostAuthLoginOK()
}
