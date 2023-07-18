package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"notes-rew/internal/middlewares"
	"notes-rew/internal/users_service/controller"
	"notes-rew/internal/users_service/models"
	"notes-rew/internal/users_service/usecase"
	"strings"
)

type UserUsecase interface {
	ReadUser(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUser(ctx context.Context, req usecase.UpdateUserInput) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
}

type UserController struct {
	usecase   UserUsecase
	validator *validator.Validate
}

func (c *UserController) Register(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Use(middlewares.SessionMiddleware)
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

	session, ok := ctx.Value("session").(*sessions.Session)
	if !ok {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	currentUserID := session.Values["userID"].(uuid.UUID)

	user, err := c.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	session, ok := ctx.Value("session").(*sessions.Session)
	if !ok {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	currentUserID, ok := session.Values["userID"].(uuid.UUID)
	if !ok {
		http.Error(w, "no userID", http.StatusInternalServerError)
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

	existingUser, err := c.usecase.CheckUserByEmail(r.Context(), strings.ToLower(*req.Email))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingUser {
		http.Error(w, "email already exists", http.StatusBadRequest)
		return
	}

	if err := c.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	session, ok := ctx.Value("session").(*sessions.Session)
	if !ok {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	currentUserID, ok := session.Values["userID"].(uuid.UUID)
	if !ok {
		http.Error(w, "no userID", http.StatusInternalServerError)
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

	session.Values["userID"] = nil
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewUserController(usecase UserUsecase, validate *validator.Validate) *UserController {
	return &UserController{
		usecase:   usecase,
		validator: validate,
	}
}
