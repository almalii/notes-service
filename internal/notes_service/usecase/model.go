package usecase

import (
	"time"

	"github.com/google/uuid"
)

type CreateNoteInput struct {
	Title  string   `json:"title" validate:"required,alphanum,min=1,max=50"`
	Body   string   `json:"body" validate:"required,bytesize"`
	Tags   []string `json:"tags" validate:"omitempty"`
	Author uuid.UUID
}

type UpdateNoteInput struct {
	Title     *string   `json:"title" validate:"required,alphanum,min=1,max=50"`
	Body      *string   `json:"body" validate:"required,bytesize"`
	Tags      *[]string `json:"tags" validate:"omitempty"`
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

type ReadNoteInput struct {
	NoteID        uuid.UUID
	CurrentUserID uuid.UUID
}
