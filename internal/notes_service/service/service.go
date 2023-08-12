package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/notes_service/models"
	"time"
)

const ExpirationTime = time.Hour * 24

type NoteStorage interface {
	CreateNoteByID(ctx context.Context, note CreateNote) error
	GetNoteByID(ctx context.Context, id uuid.UUID) (models.NoteOutput, error)
	GetAllNotesByAuthorID(ctx context.Context, currentUserID uuid.UUID) ([]models.NoteOutput, error)
	UpdateNoteByID(ctx context.Context, id uuid.UUID, note UpdateNote) error
	DeleteNoteByID(ctx context.Context, id uuid.UUID) error
}

type NoteService struct {
	storage NoteStorage
	cache   *redis.Client
}

func (s *NoteService) SaveNoteByID(ctx context.Context, note CreateNote) error {
	if err := s.storage.CreateNoteByID(ctx, note); err != nil {
		return err
	}

	noteJSON, _ := json.Marshal(note)
	if err := s.cache.Set(ctx, note.ID.String(), noteJSON, ExpirationTime).Err(); err != nil {
		return err
	}

	return nil
}

func (s *NoteService) GetNoteByID(ctx context.Context, id uuid.UUID) (*models.NoteOutput, error) {
	var note models.NoteOutput

	cachedNote, err := s.cache.Get(ctx, id.String()).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedNote), &note); err != nil {
			logrus.Printf("error while unmarshaling cached data: %v", err)
		}
		return &note, nil
	}

	// Получаем данные из базы данных
	note, err = s.storage.GetNoteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Если кеш доступен, пытаемся сохранить данные в нем
	if s.cache != nil {
		noteJSON, _ := json.Marshal(note)
		err = s.cache.Set(ctx, id.String(), noteJSON, ExpirationTime).Err()
		if err != nil {
			logrus.Printf("error while saving to Redis: %v", err)
		}
	}

	return &note, nil
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

func NewNoteService(storage NoteStorage, client *redis.Client) *NoteService {
	return &NoteService{
		storage: storage,
		cache:   client,
	}
}
