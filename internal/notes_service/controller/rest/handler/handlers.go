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
	"notes-rew/internal/token_manager"
)

type NoteUsecase interface {
	CreateNote(ctx context.Context, req usecase.CreateNoteInput) (uuid.UUID, error)
	ReadNote(ctx context.Context, noteID, currentUserID uuid.UUID) (models.NoteOutput, error)
	ReadAllNotes(ctx context.Context, currentUserID uuid.UUID) ([]models.NoteOutput, error)
	UpdateNote(ctx context.Context, id uuid.UUID, req usecase.UpdateNoteInput) error
	DeleteNote(ctx context.Context, id uuid.UUID) error
}

type NoteController struct {
	usecase      NoteUsecase
	validator    *validator.Validate
	tokenManager token_manager.TokenManager
}

func (c *NoteController) Register(r chi.Router) {
	r.Route("/notes", func(r chi.Router) {
		r.Use(middlewares.UserIdentity(c.tokenManager))
		r.Post("/", c.CreateNoteHandler)
		r.Get("/{id}", c.GetNoteHandler)
		r.Get("/", c.GetAllNotesHandler)
		r.Patch("/{id}", c.UpdateNoteHandler)
		r.Delete("/{id}", c.DeleteNoteHandler)
	})
}

// @Summary CreateNote
// @Description create note
// @Security JWTAuth
// @Tags notes
// @Accept json
// @Produce json
// @Param note body controller.CreateNoteRequest true "Note info"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /notes [post]
func (c *NoteController) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	currentUserID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	var req controller.CreateNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := c.validator.Struct(req); err != nil {
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
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary GetNote
// @Description get note
// @Security JWTAuth
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /notes/{id} [get]
func (c *NoteController) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
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

	noteID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(noteID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	note, err := c.usecase.ReadNote(ctx, parsedUUID, currentUserID)
	if err != nil {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary GetAllNotes
// @Description get all notes
// @Security JWTAuth
// @Tags notes
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /notes [get]
func (c *NoteController) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
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

	notes, err := c.usecase.ReadAllNotes(ctx, currentUserID)
	if err != nil {
		http.Error(w, "failed to retrieve notes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary UpdateNote
// @Description update note
// @Security JWTAuth
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Param note body controller.UpdateNoteRequest true "Note info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /notes/{id} [patch]
func (c *NoteController) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	currentUserID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "error reading id", http.StatusNotFound)
		return
	}

	noteID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(noteID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	_, err = c.usecase.ReadNote(ctx, parsedUUID, currentUserID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	var req controller.UpdateNoteRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Error(err)
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

	err = c.usecase.UpdateNote(ctx, parsedUUID, domain)
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

// @Summary DeleteNote
// @Description delete note
// @Security JWTAuth
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /notes/{id} [delete]
func (c *NoteController) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
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

	noteID := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(noteID)
	if err != nil {
		logrus.Error("error converting string to UUID", err)
		return
	}

	_, err = c.usecase.ReadNote(ctx, parsedUUID, currentUserID)
	if err != nil {
		http.Error(w, "id is not found", http.StatusNotFound)
		return
	}

	err = c.usecase.DeleteNote(ctx, parsedUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewNoteController(
	usecase NoteUsecase,
	validator *validator.Validate,
	tokenManager token_manager.TokenManager,
) *NoteController {
	return &NoteController{
		usecase:      usecase,
		validator:    validator,
		tokenManager: tokenManager,
	}
}
