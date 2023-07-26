package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"notes-rew/internal/sessions"
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
	usecase      UserUsecase
	validator    *validator.Validate
	sessionStore *sessions.SessionStore
}

func (c *UserController) Register(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
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

	currentSessionID, err := sessions.GetSessionByCookie(r, "session-id")
	if err != nil {
		logrus.Error("error getting session id from cookie: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getSession, err := c.sessionStore.Get(ctx, currentSessionID)
	if err != nil {
		logrus.Error("error getting session store: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	valueUserID := getSession.Values["userID"]
	currentUserID, err := uuid.Parse(valueUserID.(string))
	if err != nil {
		logrus.Error("error parsing userID: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	currentSessionID, err := sessions.GetSessionByCookie(r, "session-id")
	if err != nil {
		logrus.Error("error getting session id from cookie: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getSession, err := c.sessionStore.Get(ctx, currentSessionID)
	if err != nil {
		logrus.Error("error getting session store: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	valueUserID := getSession.Values["userID"]
	currentUserID, err := uuid.Parse(valueUserID.(string))
	if err != nil {
		logrus.Error("error parsing user id: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = c.usecase.ReadUser(ctx, currentUserID)
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

	currentSessionID, err := sessions.GetSessionByCookie(r, "session-id")
	if err != nil {
		logrus.Error("error getting session id from cookie: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getSession, err := c.sessionStore.Get(ctx, currentSessionID)
	if err != nil {
		logrus.Error("error getting session store: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	valueUserID := getSession.Values["userID"]
	currentUserID, err := uuid.Parse(valueUserID.(string))
	if err != nil {
		logrus.Error("error parsing user id: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = c.usecase.ReadUser(ctx, currentUserID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	err = c.usecase.DeleteUser(ctx, currentUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessions.ClearCookie(w, "session-id")
	err = c.sessionStore.Delete(ctx, currentSessionID)
	if err != nil {
		logrus.Error("error deleting session store: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewUserController(usecase UserUsecase, validate *validator.Validate, sessionStore *sessions.SessionStore) *UserController {
	return &UserController{
		usecase:      usecase,
		validator:    validate,
		sessionStore: sessionStore,
	}
}
