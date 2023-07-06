package models

import (
	"time"

	"github.com/google/uuid"
)

type NoteOutput struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	Author    uuid.UUID `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewNoteOutput(
	ID uuid.UUID,
	title string,
	body string,
	tags []string,
	author uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
) NoteOutput {
	return NoteOutput{
		ID:        ID,
		Title:     title,
		Body:      body,
		Tags:      tags,
		Author:    author,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
