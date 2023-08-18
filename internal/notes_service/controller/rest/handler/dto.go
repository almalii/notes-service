package handler

import (
	"github.com/google/uuid"
	"notes-rew/internal/notes_service/usecase"
	"time"
)

type CreateNoteRequest struct {
	Title string   `json:"title" validate:"required,alphanum,min=1,max=50"`
	Body  string   `json:"body" validate:"required,bytesize"`
	Tags  []string `json:"tags" validate:"omitempty"`
}

func (cnr CreateNoteRequest) ToDomain(uuid uuid.UUID) usecase.CreateNoteInput {
	return usecase.CreateNoteInput{
		Title:  cnr.Title,
		Body:   cnr.Body,
		Tags:   cnr.Tags,
		Author: uuid,
	}
}

type UpdateNoteRequest struct {
	Title     string    `json:"title" validate:"required,alphanum,min=1,max=50"`
	Body      string    `json:"body" validate:"required,bytesize"`
	Tags      []string  `json:"tags" validate:"omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (unr UpdateNoteRequest) ToDomain() usecase.UpdateNoteInput {
	return usecase.UpdateNoteInput{
		Title: &unr.Title,
		Body:  &unr.Body,
		Tags:  &unr.Tags,
	}
}

type NoteResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Body  string    `json:"body"`
}

func NewNoteResponse(ID uuid.UUID, title string, body string) *NoteResponse {
	return &NoteResponse{ID: ID, Title: title, Body: body}
}
