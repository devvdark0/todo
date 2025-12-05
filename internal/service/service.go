package service

import (
	"github.com/devvdark0/todo/internal/model"
	"github.com/devvdark0/todo/internal/storage"
	"github.com/google/uuid"
)

type Storage interface {
	Create(model.Task) error
	Get(id uuid.UUID) (model.Task, error)
	Update(id uuid.UUID, task model.Task) error
	List() ([]model.Task, error)
	Delete(id uuid.UUID) error
}

type TodoService struct {
	storage Storage
}

func NewService(store *storage.TodoStore) *TodoService {
	return &TodoService{storage: store}
}
