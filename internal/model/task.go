package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	IsDone      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required, max=255"`
	Description string `json:"description"`
	IsDone      bool   `json:"is_done" validate:"required"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsDone      *bool   `json:"is_done"`
}
