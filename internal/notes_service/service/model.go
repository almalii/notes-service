package service

import (
	"time"

	"github.com/google/uuid"
)

type CreateNote struct {
	ID        uuid.UUID
	Title     string
	Body      string
	Tags      []string
	Author    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCreateNote(
	ID uuid.UUID,
	title string,
	body string,
	tags []string,
	author uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
) CreateNote {
	return CreateNote{
		ID:        ID,
		Title:     title,
		Body:      body,
		Tags:      tags,
		Author:    author,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type UpdateNote struct {
	Title     *string
	Body      *string
	Tags      *[]string
	UpdatedAt time.Time
}
