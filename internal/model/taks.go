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
