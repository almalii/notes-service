package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"notes-rew/internal/middlewares"
	"notes-rew/internal/token_manager"
	"notes-rew/internal/users_service/models"
	"notes-rew/internal/users_service/usecase"
)

const userIDKey = "userID"

type UserUsecase interface {
	ReadUser(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUser(ctx context.Context, req usecase.UpdateUserInput) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserController struct {
	usecase      UserUsecase
	tokenManager *token_manager.TokenManager
	validator    *validator.Validate
}

func (c *UserController) Register(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Use(middlewares.UserIdentity(c.tokenManager))
		r.Get("/", c.GetUserHandler)
		r.Patch("/", c.UpdateUserHandler)
		r.Delete("/", c.DeleteUserHandler)
	})

}

// GetUserHandler
// @Summary GetUser
// @Description get user
// @Security JWTAuth
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /users [get]
func (c *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	currentUserID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		logrus.Error("error reading id from context")
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	user, err := c.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error reading user", err)
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Authorization", r.Header.Get("Authorization"))
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateUserHandler
// @Summary UpdateUser
// @Description update user
// @Security JWTAuth
// @Tags users
// @Accept json
// @Produce json
// @Param user body controller.UpdateUserRequest true "User info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /users [patch]
func (c *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	currentUserID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		logrus.Error("error reading id from context")
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	_, err := c.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error reading user", err)
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	var req UpdateUserRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domain := req.ToDomain(currentUserID)

	err = c.usecase.UpdateUser(ctx, domain)
	if err != nil {
		logrus.Error("error updating user", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteUserHandler
// @Summary DeleteUser
// @Description delete user
// @Security JWTAuth
// @Tags users
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 500
// @Router /users [delete]
func (c *UserController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	currentUserID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		logrus.Error("error reading id from context")
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	_, err := c.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error reading user", err)
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	err = c.usecase.DeleteUser(ctx, currentUserID)
	if err != nil {
		logrus.Error("error deleting user", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewUserController(
	usecase UserUsecase,
	tokenManager *token_manager.TokenManager,
	validator *validator.Validate,
) *UserController {
	return &UserController{
		usecase:      usecase,
		tokenManager: tokenManager,
		validator:    validator,
	}
}
