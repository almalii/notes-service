package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"notes-rew/internal/middlewares"
	"notes-rew/internal/token_manager"
	"notes-rew/internal/users_service/controller"
	"notes-rew/internal/users_service/models"
	"notes-rew/internal/users_service/usecase"
)

type UserUsecase interface {
	ReadUser(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUser(ctx context.Context, req usecase.UpdateUserInput) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserController struct {
	usecase      UserUsecase
	tokenManager token_manager.TokenManager
}

func (c *UserController) Register(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Use(middlewares.UserIdentity(c.tokenManager))
		r.Get("/", c.GetUserHandler)
		r.Patch("/", c.UpdateUserHandler)
		r.Delete("/", c.DeleteUserHandler)
	})

}

// @Summary GetUser
// @Description get user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.UserOutput
// @Failure 400 {object} integer
// @Failure 500 {object} integer
// @Router /users [get]
func (c *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	currentUserID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	user, err := c.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
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

// @Summary UpdateUser
// @Description update user
// @Tags users
// @Accept json
// @Produce json
// @Param user body controller.UpdateUserRequest true "User info"
// @Success 200 {object} controller.UpdateUserRequest
// @Failure 400 {object} integer
// @Failure 500 {object} integer
// @Router /users [patch]
func (c *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	currentUserID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	_, err := c.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	var req controller.UpdateUserRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	domain := req.ToDomain(currentUserID)

	err = c.usecase.UpdateUser(ctx, domain)
	if err != nil {
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

// @Summary DeleteUser
// @Description delete user
// @Tags users
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} integer
// @Failure 500 {object} integer
// @Router /users [delete]
func (c *UserController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	currentUserID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	_, err := c.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	err = c.usecase.DeleteUser(ctx, currentUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewUserController(usecase UserUsecase, tokenManager token_manager.TokenManager) *UserController {
	return &UserController{
		usecase:      usecase,
		tokenManager: tokenManager,
	}
}
