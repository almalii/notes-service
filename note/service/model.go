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

type UpdateNote struct {
	Title     *string
	Body      *string
	Tags      *[]string
	UpdatedAt time.Time
}
