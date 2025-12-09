package db

import (
	"database/sql"
	"github.com/devvdark0/todo/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.DbConfig.DSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
