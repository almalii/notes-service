package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"notes-rew/internal/auth_service/controller"
	"notes-rew/internal/auth_service/models"
	"notes-rew/internal/auth_service/usecase"
)

type AuthUsecase interface {
	CreateUser(ctx context.Context, req usecase.UserInput) (uuid.UUID, error)
	AuthenticateUser(ctx context.Context, req usecase.AuthInput) (*models.AuthResponse, error)
}

type AuthController struct {
	usecase AuthUsecase
}

func (c *AuthController) Register(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.SignUpHandler)
		r.Post("/login", c.SignInHandler)
		r.Post("/logout", c.SignOutHandler)
	})
}

// Регистрация пользователя
func (c *AuthController) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req controller.SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	domain := req.ToDomain()

	userID, err := c.usecase.CreateUser(ctx, domain)
	if err != nil {
		logrus.Errorf("error creating user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := controller.NewSignUpResponse(userID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
}

// Авторизация пользователя
func (c *AuthController) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req controller.SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.Errorf("error decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	domain := req.ToDomain()

	resp, err := c.usecase.AuthenticateUser(ctx, domain)
	if err != nil {
		http.Error(w, "password is not correct", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)

}

// Выход пользователя
func (c *AuthController) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewAuthController(
	usecase AuthUsecase,
) *AuthController {
	return &AuthController{
		usecase: usecase,
	}
}
