package service

import (
	"fmt"
	"github.com/devvdark0/todo/internal/storage/task"
	"time"

	"github.com/devvdark0/todo/internal/model"
	"github.com/go-playground/validator/v10"
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

func NewService(store *task.TodoStore) *TodoService {
	return &TodoService{storage: store}
}

func (s *TodoService) ListTasks() ([]model.Task, error) {
	tasks, err := s.storage.List()
	if err != nil {
		return nil, fmt.Errorf("list tasks service: %w", err)
	}

	return tasks, nil
}

func (s *TodoService) GetTaskByID(id string) (model.Task, error) {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return model.Task{}, fmt.Errorf("get task service: %w", err)
	}

	task, err := s.storage.Get(uuidId)
	if err != nil {
		return model.Task{}, fmt.Errorf("get task service: %w", err)
	}

	return task, nil
}

func (s *TodoService) CreateTask(req model.CreateTaskRequest) error {
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return fmt.Errorf("create user service: %w", err)
	}

	id := uuid.New()
	task := model.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		IsDone:      req.IsDone,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.storage.Create(task); err != nil {
		return fmt.Errorf("create user service: %w", err)
	}

	return nil
}

func (s *TodoService) UpdateTask(id string, req model.UpdateTaskRequest) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("delete task service: %w", err)
	}

	task, err := s.storage.Get(uuidId)
	if err != nil {
		return fmt.Errorf("delete task service: %w", err)
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.IsDone != nil {
		task.IsDone = *req.IsDone
	}

	if err := s.storage.Update(uuidId, task); err != nil {
		return fmt.Errorf("delete task service: %w", err)
	}

	return nil
}

func (s *TodoService) DeleteTask(id string) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("update task service: %w", err)
	}

	if err = s.storage.Delete(uuidId); err != nil {
		return fmt.Errorf("update task service: %w", err)
	}

	return nil
}
