package storage

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/note/models"
	"notes-rew/note/service"

	"github.com/jackc/pgx/v5"
)

type NoteStorage struct {
	db *pgx.Conn
}

func (s *NoteStorage) CreateNoteByID(ctx context.Context, note service.CreateNote) error {
	sql, args, err := squirrel.Insert("notes").
		Columns("id", "title", "body", "tags", "author", "created_at", "updated_at").
		Values(note.ID, note.Title, note.Body, note.Tags, note.Author, note.CreatedAt, note.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
	}

	_, err = s.db.Exec(ctx, sql, args...)
	if err != nil {
		logrus.Errorf("error while executing squirrel query: %v", err)
	}

	return nil
}

func (s *NoteStorage) GetNoteByID(ctx context.Context, id uuid.UUID) (models.NoteOutput, error) {
	var note NoteResponse

	sql, args, err := squirrel.Select("id", "title", "body", "tags", "author", "created_at", "updated_at").
		From("notes").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&note.ID, &note.Title, &note.Body, &note.Tags, &note.Author, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return models.NoteOutput{}, err
	}

	return models.NoteOutput(note), nil

}

func (s *NoteStorage) GetAllNotesByAuthorID(ctx context.Context, authorID uuid.UUID) ([]models.NoteOutput, error) {
	sql, args, err := squirrel.Select("id", "title", "body", "tags", "author", "created_at", "updated_at").
		From("notes").
		Where(squirrel.Eq{"author": authorID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.NoteOutput
	for rows.Next() {
		var note models.NoteOutput
		err := rows.Scan(&note.ID, &note.Title, &note.Body, &note.Tags, &note.Author, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *NoteStorage) UpdateNoteByID(ctx context.Context, id uuid.UUID, note service.UpdateNote) error {

	sql, args, err := squirrel.Update("notes").
		Set("title", note.Title).
		Set("body", note.Body).
		Set("tags", note.Tags).
		Set("updated_at", note.UpdatedAt).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
	}

	_, err = s.db.Exec(ctx, sql, args...)
	if err != nil {
		logrus.Errorf("error while executing squirrel query: %v", err)
	}

	return nil
}

func (s *NoteStorage) DeleteNoteByID(ctx context.Context, id uuid.UUID) error {
	sql, args, err := squirrel.Delete("notes").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
	}

	_, err = s.db.Exec(ctx, sql, args...)
	if err != nil {
		logrus.Errorf("error while executing squirrel query: %v", err)
	}

	return nil
}

func NewNoteStorage(db *pgx.Conn) *NoteStorage {
	return &NoteStorage{
		db: db,
	}
}
