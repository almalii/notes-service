package controller

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"notes-rew/note/usecase"
)

type CreateNoteRequest struct {
	Title  string   `json:"title" validate:"required"`
	Body   string   `json:"body" validate:"required"`
	Tags   []string `json:"tags" validate:",omitempty"`
	Author uuid.UUID
}

func (cnr CreateNoteRequest) ToDomain() (usecase.CreateNoteInput, error) {
	if len(cnr.Body) > 30<<20 {
		return usecase.CreateNoteInput{}, errors.New("too large note")
	}

	input, err := usecase.NewCreateNoteInput(cnr.Title, cnr.Body, cnr.Tags, cnr.Author)
	if err != nil {
		return usecase.CreateNoteInput{}, err
	}

	return input, nil
}

type UpdateNoteRequest struct {
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (unr UpdateNoteRequest) ToDomain() (usecase.UpdateNoteInput, error) {
	if len(unr.Body) > 30<<20 {
		return usecase.UpdateNoteInput{}, errors.New("too large note")
	}

	input, err := usecase.NewUpdateNoteInput(&unr.Title, &unr.Body, &unr.Tags)
	if err != nil {
		return usecase.UpdateNoteInput{}, err
	}

	return input, nil
}

type NoteResponseId struct {
	ID uuid.UUID `json:"id"`
}
