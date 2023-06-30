package handler

import (
	"bookmarks/note/controller"
	"bookmarks/note/models"
	"bookmarks/note/usecase"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type NoteUsecase interface {
	CreateNote(ctx context.Context, req usecase.CreateNoteInput, currentUserID uuid.UUID) (uuid.UUID, error)
	ReadNote(ctx context.Context, id uuid.UUID) (models.NoteOutput, error)
	ReadAllNotes(ctx context.Context, currentUserID uuid.UUID) ([]models.NoteOutput, error)
	UpdateNote(ctx context.Context, id uuid.UUID, req usecase.UpdateNoteInput) error
	DeleteNote(ctx context.Context, id uuid.UUID) error
}

type NoteController struct {
	usecase NoteUsecase
}

func (c *NoteController) Register(r chi.Router) {
	r.Route("/note", func(r chi.Router) {
		r.Use(SessionMiddleware)
		r.Post("/", c.CreateNoteHandler)
		r.Get("/{id}", c.GetNoteHandler)
		r.Get("/all", c.GetAllNotesHandler)
		r.Put("/{id}", c.UpdateNoteHandler)
		r.Delete("/{id}", c.DeleteNoteHandler)
	})
}

func (c *NoteController) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)

	var req controller.CreateNoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	domain, err := req.ToDomain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUserID := session.Values["userID"].(uuid.UUID)

	noteID, err := c.usecase.CreateNote(r.Context(), domain, currentUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := controller.NoteResponseId{
		ID: noteID,
	}

	WriteJSONResponse(w, http.StatusCreated, resp)

}

func (c *NoteController) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)
	currentUserID := session.Values["userID"].(uuid.UUID)

	noteID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(noteID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	note, err := c.usecase.ReadNote(r.Context(), parsedUUID)
	if err != nil {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	if note.Author != currentUserID {
		http.Error(w, "not authorized to read this note", http.StatusUnauthorized)
		return
	}

	WriteJSONResponse(w, http.StatusOK, note)
}

func (c *NoteController) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)
	currentUserID := session.Values["userID"].(uuid.UUID)

	notes, err := c.usecase.ReadAllNotes(r.Context(), currentUserID)
	if err != nil {
		http.Error(w, "failed to retrieve notes", http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, http.StatusOK, notes)
}

func (c *NoteController) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)
	currentUserID := session.Values["userID"].(uuid.UUID)

	noteID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(noteID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	note, err := c.usecase.ReadNote(r.Context(), parsedUUID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	if note.Author != currentUserID {
		http.Error(w, "not authorized to update this note", http.StatusUnauthorized)
		return
	}

	var req controller.UpdateNoteRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	domain, err := req.ToDomain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.UpdatedAt = time.Now()

	err = c.usecase.UpdateNote(r.Context(), parsedUUID, domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, http.StatusOK, req)
}

func (c *NoteController) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)
	currentUserID := session.Values["userID"].(uuid.UUID)

	noteID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(noteID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	note, err := c.usecase.ReadNote(r.Context(), parsedUUID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	if note.Author != currentUserID {
		http.Error(w, "not authorized to delete this note", http.StatusUnauthorized)
		return
	}

	err = c.usecase.DeleteNote(r.Context(), parsedUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewNoteController(usecase NoteUsecase) *NoteController {
	return &NoteController{
		usecase: usecase,
	}
}
