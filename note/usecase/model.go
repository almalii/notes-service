package usecase

import (
	"github.com/google/uuid"
	"time"
)

type CreateNoteInput struct {
	ID        uuid.UUID
	Title     string
	Body      string
	Tags      []string
	Author    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCreateNoteInput(title string, body string, tags []string, currentUserID uuid.UUID) (CreateNoteInput, error) {
	return CreateNoteInput{
		ID:        uuid.New(),
		Title:     title,
		Body:      body,
		Tags:      tags,
		Author:    currentUserID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

type UpdateNoteInput struct {
	Title     *string
	Body      *string
	Tags      *[]string
	UpdatedAt time.Time
}

func NewUpdateNoteInput(title *string, body *string, tags *[]string) (UpdateNoteInput, error) {
	return UpdateNoteInput{
		Title:     title,
		Body:      body,
		Tags:      tags,
		UpdatedAt: time.Now().UTC(),
	}, nil
}
