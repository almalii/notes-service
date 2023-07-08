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
	"notes-rew/internal/auth/controller"
	"notes-rew/internal/auth/models"
	"notes-rew/internal/auth/usecase"
	"notes-rew/internal/middlewares"
	"strings"
)

type AuthUsecase interface {
	CreateUser(ctx context.Context, req usecase.UserInput) (uuid.UUID, error)
	AuthenticateUser(ctx context.Context, req usecase.AuthInput) (models.AuthResponse, error)
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
}

type AuthController struct {
	usecase   AuthUsecase
	validator *validator.Validate
}

func (c *AuthController) Register(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Use(middlewares.SessionMiddleware)
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

	session, ok := r.Context().Value("session").(*sessions.Session)
	if !ok {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	var req controller.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	existingUser, err := c.usecase.CheckUserByEmail(r.Context(), strings.ToLower(req.Email))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingUser {
		http.Error(w, "email already exists", http.StatusBadRequest)
		return
	}

	domain, err := req.ToDomain(c.validator)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := c.usecase.CreateUser(r.Context(), domain)
	if err != nil {
		logrus.Errorf("error creating user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = userID
	err = session.Save(r, w)
	if err != nil {
		logrus.Errorf("error saving session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	session, ok := r.Context().Value("session").(*sessions.Session)
	if !ok {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	if session.Values["userID"] != nil {
		http.Error(w, "already signed in", http.StatusBadRequest)
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

	domain, err := req.ToDomain(c.validator)
	if err != nil {
		logrus.Errorf("error converting to domain: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	resp, err := c.usecase.AuthenticateUser(r.Context(), domain)
	if err != nil {
		http.Error(w, "password is not correct", http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = resp.ID
	err = session.Save(r, w)
	if err != nil {
		logrus.Errorf("error saving session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
