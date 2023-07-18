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
	"notes-rew/internal/auth_service/controller"
	"notes-rew/internal/auth_service/models"
	"notes-rew/internal/auth_service/usecase"
	"notes-rew/internal/session"
	"strings"
	"time"
)

type AuthUsecase interface {
	CreateUser(ctx context.Context, req usecase.UserInput) (uuid.UUID, error)
	AuthenticateUser(ctx context.Context, req usecase.AuthInput) (models.AuthResponse, error)
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
}

type AuthController struct {
	usecase      AuthUsecase
	validator    *validator.Validate
	sessionStore *session.RedisSessionStore
}

func (c *AuthController) Register(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		//r.Use(middlewares.SessMiddleware)
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

	sessionStore := session.NewRedisSessionStore("0.0.0.0:32768", "", 1)
	c.sessionStore = sessionStore

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

	sess, err := session.CreateSession(12 * time.Hour)
	sess.Values["id"] = userID
	err = sessionStore.Set(ctx, sess)
	if err != nil {
		logrus.Errorf("Failed to save session in Redis: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := controller.NewSignUpResponse(userID, domain.Username, domain.Email)

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

	sessionStore := session.NewRedisSessionStore("0.0.0.0:32768", "", 1)
	sessionID := r.Header.Get("X-Session-ID")
	sessions, err := sessionStore.Get(ctx, sessionID)
	if err != nil {
		logrus.Errorf("Failed to get session from Redis: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req controller.SignInRequest
	err = json.NewDecoder(r.Body).Decode(&req)
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

	existingSession, err := sessionStore.Get(ctx, resp.ID.String())
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existingSession != nil {
		http.Error(w, "user is already signed in", http.StatusBadRequest)
		return
	}

	session, err := session.CreateSession(1 * time.Minute)
	session.Values["id"] = resp.ID
	err = sessionStore.Set(ctx, sessions)
	if err != nil {
		logrus.Errorf("error creating session: %v", err)
		return
	}

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

	session, ok := r.Context().Value("session").(*sessions.Session)
	if !ok {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = nil
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		logrus.Errorf("error saving session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewAuthController(usecase AuthUsecase, validator *validator.Validate) *AuthController {
	return &AuthController{
		usecase:   usecase,
		validator: validator,
	}
}
