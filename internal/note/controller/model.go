package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/note/usecase"
	"time"
)

type CreateNoteRequest struct {
	Title string   `json:"title" validate:"required,alphanum,min=1,max=50"`
	Body  string   `json:"body" validate:"required,bytesize"`
	Tags  []string `json:"tags" validate:"omitempty"`
}

func (cnr CreateNoteRequest) ToDomain(uuid uuid.UUID, validate *validator.Validate) (usecase.CreateNoteInput, error) {
	if err := validate.Struct(cnr); err != nil {
		logrus.Error(err)
		return usecase.CreateNoteInput{}, err.(validator.ValidationErrors)
	}

	return usecase.CreateNoteInput{
		Title:  cnr.Title,
		Body:   cnr.Body,
		Tags:   cnr.Tags,
		Author: uuid,
	}, nil
}

type UpdateNoteRequest struct {
	Title     string    `json:"title" validate:"required,alphanum,min=1,max=50"`
	Body      string    `json:"body" validate:"required,bytesize"`
	Tags      []string  `json:"tags" validate:"omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (unr UpdateNoteRequest) ToDomain(validate *validator.Validate) (usecase.UpdateNoteInput, error) {
	if err := validate.Struct(unr); err != nil {
		logrus.Error(err)
		return usecase.UpdateNoteInput{}, err.(validator.ValidationErrors)
	}

	return usecase.UpdateNoteInput{
		Title: &unr.Title,
		Body:  &unr.Body,
		Tags:  &unr.Tags,
	}, nil
}

type NoteResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Body  string    `json:"body"`
}
