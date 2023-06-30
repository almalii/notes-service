package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"notes-rew/user/controller"
	"notes-rew/user/models"
	"notes-rew/user/usecase"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req usecase.CreateUserInput) (uuid.UUID, error)
	ReadUser(ctx context.Context, id uuid.UUID) (models.UserOutput, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req usecase.UpdateUserInput) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	AuthenticateUser(ctx context.Context, req usecase.AuthInput) (models.AuthResponse, error)
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
}

type UserController struct {
	usecase UserUsecase
}

func (c *UserController) Register(r chi.Router) {

	r.Route("/user", func(r chi.Router) {
		r.Use(SessionMiddleware)
		r.Post("/", c.CreateUserHandler)
		r.Get("/{id}", c.GetUserHandler)
		r.Put("/", c.UpdateUserHandler)
		r.Delete("/", c.DeleteUserHandler)
	})

	r.Route("/login", func(r chi.Router) {
		r.Use(SessionMiddleware)
		r.Post("/", c.SignInHandler)
	})

	r.Route("/logout", func(r chi.Router) {
		r.Use(SessionMiddleware)
		r.Post("/", c.SignOutHandler)
	})

}

// TODO SignUpHandler
func (c *UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)

	req := new(controller.CreateUserRequest)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	existingUser, err := c.usecase.CheckUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingUser {
		http.Error(w, "email already exists", http.StatusBadRequest)
		return
	}

	domain, err := req.ToDomain()
	if err != nil {
		logrus.Errorf("error converting request to domain: %s", err)
	}

	userID, err := c.usecase.CreateUser(r.Context(), domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = userID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := controller.UserResponseId{
		ID: userID,
	}

	WriteJSONResponse(w, http.StatusCreated, resp)

}

func (c *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	user, err := c.usecase.ReadUser(r.Context(), parsedUUID)
	if err != nil {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	user.UpdatedAt = user.UpdatedAt.Local()
	user.CreatedAt = user.CreatedAt.Local()

	WriteJSONResponse(w, http.StatusOK, user)
}

func (c *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get session
	session := r.Context().Value("session").(*sessions.Session)
	currentUserID := session.Values["userID"].(uuid.UUID)

	_, err := c.usecase.ReadUser(r.Context(), currentUserID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	req := new(controller.UpdateUserRequest)
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	domain, err := req.ToDomain()
	if err != nil {
		logrus.Errorf("error converting to domain: %v", err)
	}

	err = c.usecase.UpdateUser(r.Context(), currentUserID, domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, http.StatusOK, req)
}

func (c *UserController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)
	currentUserID := session.Values["userID"].(uuid.UUID)

	_, err := c.usecase.ReadUser(r.Context(), currentUserID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	err = c.usecase.DeleteUser(r.Context(), currentUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *UserController) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)
	if session.Values["userID"] != nil {
		http.Error(w, "already signed in", http.StatusBadRequest)
		return
	}

	req := new(controller.AuthRequest)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	domain, err := req.ToDomain()
	if err != nil {
		logrus.Errorf("error converting to domain: %v", err)
	}

	resp, err := c.usecase.AuthenticateUser(r.Context(), domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = resp.ID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, http.StatusOK, resp)

}

func (c *UserController) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)

	session.Values["userID"] = nil
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewUserController(usecase UserUsecase) *UserController {
	return &UserController{
		usecase: usecase,
	}
}
