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
	"notes-rew/internal/notes_service/controller"
	"notes-rew/internal/notes_service/models"
	"notes-rew/internal/notes_service/usecase"
	"notes-rew/internal/sessions"
)

type NoteUsecase interface {
	CreateNote(ctx context.Context, req usecase.CreateNoteInput) (uuid.UUID, error)
	ReadNote(ctx context.Context, id uuid.UUID) (models.NoteOutput, error)
	ReadAllNotes(ctx context.Context, currentUserID uuid.UUID) ([]models.NoteOutput, error)
	UpdateNote(ctx context.Context, id uuid.UUID, req usecase.UpdateNoteInput) error
	DeleteNote(ctx context.Context, id uuid.UUID) error
}

type NoteController struct {
	usecase      NoteUsecase
	validator    *validator.Validate
	sessionStore *sessions.SessionStore
}

func (c *NoteController) Register(r chi.Router) {
	r.Route("/notes", func(r chi.Router) {
		r.Use(middlewares.SessionMiddleware)
		r.Post("/", c.CreateNoteHandler)
		r.Get("/{id}", c.GetNoteHandler)
		r.Get("/", c.GetAllNotesHandler)
		r.Put("/{id}", c.UpdateNoteHandler)
		r.Delete("/{id}", c.DeleteNoteHandler)
	})
}

func (c *NoteController) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
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

	var req controller.CreateNoteRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = c.validator.Struct(req); err != nil {
		logrus.Error(err.(validator.ValidationErrors))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domain := req.ToDomain(currentUserID)

	noteID, err := c.usecase.CreateNote(ctx, domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := controller.NoteResponse{
		ID:    noteID,
		Title: domain.Title,
		Body:  domain.Body,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c *NoteController) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
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

	noteID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(noteID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	note, err := c.usecase.ReadNote(ctx, parsedUUID)
	if err != nil {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	if note.Author != currentUserID {
		http.Error(w, "not authorized to read this notes_service", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c *NoteController) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
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

	notes, err := c.usecase.ReadAllNotes(ctx, currentUserID)
	if err != nil {
		http.Error(w, "failed to retrieve notes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *NoteController) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
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
		logrus.Error("error parsing userID: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	noteID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(noteID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	note, err := c.usecase.ReadNote(ctx, parsedUUID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	if note.Author != currentUserID {
		http.Error(w, "not authorized to update this notes_service", http.StatusUnauthorized)
		return
	}

	var req controller.UpdateNoteRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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

	err = c.usecase.UpdateNote(ctx, parsedUUID, domain)
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

func (c *NoteController) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
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
		logrus.Error("error parsing userID: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	noteID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(noteID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	note, err := c.usecase.ReadNote(ctx, parsedUUID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	if note.Author != currentUserID {
		http.Error(w, "not authorized to delete this notes_service", http.StatusUnauthorized)
		return
	}

	err = c.usecase.DeleteNote(ctx, parsedUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewNoteController(usecase NoteUsecase, validator *validator.Validate, sessionStore *sessions.SessionStore) *NoteController {
	return &NoteController{
		usecase:      usecase,
		validator:    validator,
		sessionStore: sessionStore,
	}
}
