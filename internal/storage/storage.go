package storage

import (
	"database/sql"

	"github.com/devvdark0/todo/internal/model"
	"github.com/google/uuid"
)

type TodoStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *TodoStore {
	return &TodoStore{db: db}
}

func (s *TodoStore) Create(task model.Task) error {
	query := `INSERT INTO tasks (id, title, description, is_done, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, task.ID, task.Title, task.Description, task.IsDone, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoStore) Get(id uuid.UUID) (model.Task, error) {

}

func (s *TodoStore) Update(id uuid.UUID, task model.Task) error {

}

func (s *TodoStore) List() ([]model.Task, error) {

}

func (s *TodoStore) Delete(id uuid.UUID) error {

}
