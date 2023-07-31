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
		r.Put("/", c.UpdateUserHandler)
		r.Delete("/", c.DeleteUserHandler)
	})

}

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
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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
	err = json.NewEncoder(w).Encode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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
