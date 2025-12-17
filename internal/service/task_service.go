package service

import (
	"fmt"
	"time"

	"github.com/devvdark0/todo/internal/model"
	"github.com/devvdark0/todo/internal/storage"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TaskStorage interface {
	Create(task model.Task) error
	GetByID(taskID, userID uuid.UUID) (*model.Task, error)
	GetByTitle(title string, userID uuid.UUID) (*model.Task, error)
	Update(taskID, userID uuid.UUID, task model.Task) error
	List(userID uuid.UUID) ([]model.Task, error)
	Delete(taskID, userID uuid.UUID) error
}

type TodoService struct {
	storage TaskStorage
}

func NewService(store *storage.TodoStore) *TodoService {
	return &TodoService{storage: store}
}

func (s *TodoService) ListTasks(userID string) ([]model.Task, error) {
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("user id parsing err: %w", err)
	}
	tasks, err := s.storage.List(uuidUserID)
	if err != nil {
		return nil, fmt.Errorf("list tasks service: %w", err)
	}

	return tasks, nil
}

func (s *TodoService) GetTaskByID(taskID, userID string) (*model.Task, error) {
	uuidTaskID, err := uuid.Parse(taskID)
	if err != nil {
		return nil, fmt.Errorf("get task service: %w", err)
	}

	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("get task service: %w", err)
	}

	task, err := s.storage.GetByID(uuidTaskID, uuidUserID)
	if err != nil {
		return nil, fmt.Errorf("get task service: %w", err)
	}

	return task, nil
}

func (s *TodoService) CreateTask(req model.CreateTaskRequest) error {
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return fmt.Errorf("create user service: %w", err)
	}

	id := uuid.New()
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return fmt.Errorf("create task service: %w", err)
	}

	task := model.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		IsDone:      req.IsDone,
		UserId:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.storage.Create(task); err != nil {
		return fmt.Errorf("create user service: %w", err)
	}

	return nil
}

func (s *TodoService) UpdateTask(taskID, userID string, req model.UpdateTaskRequest) error {
	uuidTaskID, err := uuid.Parse(taskID)
	if err != nil {
		return fmt.Errorf("update task service: %w", err)
	}

	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("update task service: %w", err)
	}

	task, err := s.storage.GetByID(uuidTaskID, uuidUserID)
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

	if err := s.storage.Update(uuidTaskID, uuidUserID, *task); err != nil {
		return fmt.Errorf("delete task service: %w", err)
	}

	return nil
}

func (s *TodoService) DeleteTask(taskID, userID string) error {
	uuidTaskID, err := uuid.Parse(taskID)
	if err != nil {
		return fmt.Errorf("delete task service: %w", err)
	}
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("delete task service: %w", err)
	}

	if err = s.storage.Delete(uuidTaskID, uuidUserID); err != nil {
		return fmt.Errorf("update task service: %w", err)
	}

	return nil
}
