package usecase

import (
	"context"
	"fmt"
	"time"

	"notes-rew/internal/notes_service/models"
	"notes-rew/internal/notes_service/service"

	"github.com/google/uuid"
)

type NoteService interface {
	SaveNoteByID(ctx context.Context, note service.CreateNote) error
	GetNoteByID(ctx context.Context, id uuid.UUID) (models.NoteOutput, error)
	GetAllNotesByAuthorID(ctx context.Context, currentUserID uuid.UUID) ([]models.NoteOutput, error)
	UpdateNoteByID(ctx context.Context, id uuid.UUID, note service.UpdateNote) error
	DeleteNoteByID(ctx context.Context, id uuid.UUID) error
}

type NoteUsecase struct {
	service NoteService
}

func (u *NoteUsecase) CreateNote(ctx context.Context, req CreateNoteInput) (uuid.UUID, error) {
	createNote := service.NewCreateNote(
		uuid.New(),
		req.Title,
		req.Body,
		req.Tags,
		req.Author,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	err := u.service.SaveNoteByID(ctx, createNote)
	if err != nil {
		return uuid.Nil, err
	}

	return createNote.ID, nil
}

func (u *NoteUsecase) ReadNote(ctx context.Context, noteID, currentUserID uuid.UUID) (models.NoteOutput, error) {
	note, err := u.service.GetNoteByID(ctx, noteID)
	if err != nil {
		return models.NoteOutput{}, err
	}

	if note.Author != currentUserID {
		return models.NoteOutput{}, fmt.Errorf("user is not author of this note")
	}

	return note, nil
}

func (u *NoteUsecase) ReadAllNotes(ctx context.Context, currentUserID uuid.UUID) ([]models.NoteOutput, error) {
	return u.service.GetAllNotesByAuthorID(ctx, currentUserID)
}

func (u *NoteUsecase) UpdateNote(ctx context.Context, id uuid.UUID, req UpdateNoteInput) error {
	noteUpdate, err := NewUpdateNoteInput(req.Title, req.Body, req.Tags)
	if err != nil {
		return err
	}

	return u.service.UpdateNoteByID(ctx, id, service.UpdateNote(noteUpdate))
}

func (u *NoteUsecase) DeleteNote(ctx context.Context, id uuid.UUID) error {
	return u.service.DeleteNoteByID(ctx, id)
}

func NewNoteUsecase(service NoteService) *NoteUsecase {
	return &NoteUsecase{
		service: service,
	}
}
