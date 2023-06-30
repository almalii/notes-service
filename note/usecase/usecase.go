package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/note/models"
	"notes-rew/note/service"
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

func (u *NoteUsecase) CreateNote(ctx context.Context, req CreateNoteInput, currentUserID uuid.UUID) (uuid.UUID, error) {
	newNote, err := NewCreateNoteInput(req.Title, req.Body, req.Tags, currentUserID)
	if err != nil {
		logrus.Errorf("error creating note: %v", err)
	}

	err = u.service.SaveNoteByID(ctx, service.CreateNote(newNote))
	if err != nil {
		logrus.Errorf("error saving note: %v", err)
	}

	return newNote.ID, nil
}

func (u *NoteUsecase) ReadNote(ctx context.Context, id uuid.UUID) (models.NoteOutput, error) {
	return u.service.GetNoteByID(ctx, id)
}

func (u *NoteUsecase) ReadAllNotes(ctx context.Context, currentUserID uuid.UUID) ([]models.NoteOutput, error) {
	return u.service.GetAllNotesByAuthorID(ctx, currentUserID)
}

func (u *NoteUsecase) UpdateNote(ctx context.Context, id uuid.UUID, req UpdateNoteInput) error {
	noteUpdate, err := NewUpdateNoteInput(req.Title, req.Body, req.Tags)
	if err != nil {
		logrus.Errorf("error updating note: %v", err)
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
