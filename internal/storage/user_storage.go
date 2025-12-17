package storage

import (
	"database/sql"

	"github.com/devvdark0/todo/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserStore struct {
	db  *sql.DB
	log *zap.Logger
}

func NewUserStore(db *sql.DB, log *zap.Logger) *UserStore {
	return &UserStore{
		db:  db,
		log: log,
	}
}

func (s *UserStore) Create(user model.User) error {
	query := `INSERT INTO users(id, email, username, password) VALUES(?,?,?,?)`
	_, err := s.db.Exec(query, user.ID, user.Email, user.Username, user.Password)
	if err != nil {
		s.log.Error("db insert user error", zap.Error(err))
		return err
	}

	return nil
}

func (s *UserStore) GetByID(id uuid.UUID) (*model.User, error) {
	query := `SELECT id, email, username, password FROM users WHERE id=?`
	var user model.User
	err := s.db.QueryRow(query, id).
		Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Password,
		)
	if err != nil {
		s.log.Error("db select user error", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) GetByEmail(email string) (*model.User, error) {
	query := `SELECT id, email, username, password FROM users WHERE email = ?`
	var user model.User
	err := s.db.QueryRow(query, email).
		Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Password,
		)
	if err != nil {
		s.log.Error("db select user error", zap.Error(err))
		return nil, err
	}

	return &user, nil
}
