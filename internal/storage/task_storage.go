package storage

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/devvdark0/todo/internal/model"
	"github.com/google/uuid"
)

type TodoStore struct {
	db  *sql.DB
	log *zap.Logger
}

func NewStore(db *sql.DB, log *zap.Logger) *TodoStore {
	return &TodoStore{db: db, log: log}
}

func (s *TodoStore) Create(task model.Task) error {
	query := `INSERT INTO tasks (id, title, description, is_done, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(
		query,
		task.ID,
		task.Title,
		task.Description,
		task.IsDone,
		task.IsDone,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		s.log.Error("db insert err", zap.Error(err))
		return err
	}
	return nil
}

func (s *TodoStore) GetByID(taskID, userID uuid.UUID) (*model.Task, error) {
	var task model.Task
	query := `SELECT id, title, description, is_done, created_at, updated_at FROM tasks WHERE id = ? AND user_id=?`
	err := s.db.QueryRow(query, taskID, userID).
		Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.IsDone,
			&task.UserId,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
	if err != nil {
		s.log.Error("db select error", zap.Error(err))
		return nil, err
	}

	return &task, nil
}

func (s *TodoStore) GetByTitle(title string, userID uuid.UUID) (*model.Task, error) {
	query := `SELECT id, title, description, is_done, user_id, created_at, updated_at FROM tasks WHERE title LIKE '%?%' AND user_id=?`
	var task model.Task

	err := s.db.QueryRow(query, title, userID).
		Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.IsDone,
			&task.UserId,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

	if err != nil {
		s.log.Error("db selecting task by title err", zap.Error(err))
		return nil, err
	}

	return &task, nil
}

func (s *TodoStore) Update(taskID, userID uuid.UUID, task model.Task) error {
	query := `UPDATE tasks SET title=?, description=?, is_done=? WHERE id=? AND user_id=?`
	_, err := s.db.Exec(query, task.Title, task.Description, task.IsDone, taskID, userID)
	if err != nil {
		s.log.Error("db update error", zap.Error(err), zap.String("task_id", taskID.String()))
		return err
	}

	return nil
}

func (s *TodoStore) List(userID uuid.UUID) ([]model.Task, error) {
	tasks := make([]model.Task, 0)
	query := `SELECT id, title, description, is_done, user_id, created_at, updated_at FROM tasks WHERE user_id=?`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		s.log.Error("db select err", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var task model.Task
		rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.IsDone,
			&task.UserId,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *TodoStore) Delete(taskID, userID uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id=? AND user_id=?`
	_, err := s.db.Exec(query, taskID, userID)
	if err != nil {
		s.log.Error("db delete error", zap.Error(err), zap.String("id", taskID.String()))
		return err
	}

	return nil
}
