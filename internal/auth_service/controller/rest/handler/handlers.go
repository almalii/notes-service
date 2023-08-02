package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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
	usecase   AuthUsecase
	validator *validator.Validate
}

func (c *AuthController) Register(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.SignUpHandler)
		r.Post("/login", c.SignInHandler)
	})
}

// @Summary SignUp
// @Description create user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body controller.SignUpRequest true "User info"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /auth/register [post]
func (c *AuthController) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req controller.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := c.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary SignIn
// @Description login user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body controller.SignInRequest true "User info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /auth/login [post]
func (c *AuthController) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req controller.SignInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := c.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domain := req.ToDomain()

	resp, err := c.usecase.AuthenticateUser(ctx, domain)
	if err != nil {
		http.Error(w, "password is not correct", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewAuthController(usecase AuthUsecase, validator *validator.Validate) *AuthController {
	return &AuthController{
		usecase:   usecase,
		validator: validator,
	}
}
