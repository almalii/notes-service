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
	"notes-rew/internal/sessions"
	"strings"
)

type AuthUsecase interface {
	CreateUser(ctx context.Context, req usecase.UserInput) (uuid.UUID, error)
	AuthenticateUser(ctx context.Context, req usecase.AuthInput) (models.AuthResponse, error)
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
}

type AuthController struct {
	usecase      AuthUsecase
	validator    *validator.Validate
	sessionStore *sessions.SessionStore
}

func (c *AuthController) Register(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		//r.Use(middlewares.SessionMiddleware)
		r.Post("/register", c.SignUpHandler)
		r.Post("/login", c.SignInHandler)
		r.Post("/logout", c.SignOutHandler)
	})
}

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

	existingUser, err := c.usecase.CheckUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingUser {
		http.Error(w, "email already exists", http.StatusBadRequest)
		return
	}

	if err = c.validator.Struct(req); err != nil {
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

	session := sessions.NewSession()
	session.Values["userID"] = userID
	if err = c.sessionStore.Save(ctx, session); err != nil {
		logrus.Errorf("error save session in redis: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessions.SetCookie(w, session.ID, "session-id")

	resp := controller.NewSignUpResponse(userID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logrus.Errorf("error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *AuthController) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	if sessions.CheckCookieValue(w, r, "session-id") {
		http.Error(w, "user already logged in", http.StatusBadRequest)
		return
	}

	var req controller.SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.Errorf("error decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = c.validator.Struct(req); err != nil {
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

	session := sessions.NewSession()
	session.Values["userID"] = resp.UserID
	if err = c.sessionStore.Save(ctx, session); err != nil {
		logrus.Errorf("error save session in redis: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessions.SetCookie(w, session.ID, "session-id")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c *AuthController) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	sessions.ClearCookie(w, "session-id")
	err = c.sessionStore.Delete(ctx, currentSessionID)
	if err != nil {
		logrus.Error("error deleting session store: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewAuthController(usecase AuthUsecase, validator *validator.Validate, sessionStore *sessions.SessionStore) *AuthController {
	return &AuthController{
		usecase:      usecase,
		validator:    validator,
		sessionStore: sessionStore,
	}
}
