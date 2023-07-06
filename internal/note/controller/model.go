package controller

import (
	"errors"
	"notes-rew/internal/note/usecase"
	"time"

	"github.com/google/uuid"
)

type CreateNoteRequest struct {
	Title string   `json:"title" validate:"required"`
	Body  string   `json:"body" validate:"required"`
	Tags  []string `json:"tags" validate:",omitempty"`
}

func (cnr CreateNoteRequest) ToDomain(uuid uuid.UUID) (usecase.CreateNoteInput, error) {
	if len(cnr.Body) > 30<<20 {
		return usecase.CreateNoteInput{}, errors.New("too large note")
	}

	input, err := usecase.NewCreateNoteInput(cnr.Title, cnr.Body, cnr.Tags, uuid)
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
