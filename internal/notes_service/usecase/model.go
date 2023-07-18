package usecase

import (
	"time"

	"github.com/google/uuid"
)

type CreateNoteInput struct {
	Title  string
	Body   string
	Tags   []string
	Author uuid.UUID
}

func NewCreateNoteInput(title string, body string, tags []string, currentUserID uuid.UUID) (CreateNoteInput, error) {
	return CreateNoteInput{
		Title:  title,
		Body:   body,
		Tags:   tags,
		Author: currentUserID,
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
