package service

import (
	"context"
	"github.com/google/uuid"
	"notes-rew/note/models"
)

type NoteStorage interface {
	CreateNoteByID(ctx context.Context, note CreateNote) error
	GetNoteByID(ctx context.Context, id uuid.UUID) (models.NoteOutput, error)
	GetAllNotesByAuthorID(ctx context.Context, currentUserID uuid.UUID) ([]models.NoteOutput, error)
	UpdateNoteByID(ctx context.Context, id uuid.UUID, note UpdateNote) error
	DeleteNoteByID(ctx context.Context, id uuid.UUID) error
}

type NoteService struct {
	storage NoteStorage
}

func (s *NoteService) SaveNoteByID(ctx context.Context, note CreateNote) error {
	return s.storage.CreateNoteByID(ctx, note)
}

func (s *NoteService) GetNoteByID(ctx context.Context, id uuid.UUID) (models.NoteOutput, error) {
	return s.storage.GetNoteByID(ctx, id)
}

func (s *NoteService) GetAllNotesByAuthorID(ctx context.Context, authorID uuid.UUID) ([]models.NoteOutput, error) {
	return s.storage.GetAllNotesByAuthorID(ctx, authorID)
}

func (s *NoteService) UpdateNoteByID(ctx context.Context, id uuid.UUID, note UpdateNote) error {
	return s.storage.UpdateNoteByID(ctx, id, note)
}

func (s *NoteService) DeleteNoteByID(ctx context.Context, id uuid.UUID) error {
	return s.storage.DeleteNoteByID(ctx, id)
}

func NewNoteService(storage NoteStorage) *NoteService {
	return &NoteService{
		storage: storage,
	}
}
