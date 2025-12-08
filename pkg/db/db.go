package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	//TODO: change the dsn to use the environment variables
	dsn := "todo-user:todo-password@tcp(localhost:3306)/todo-db?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
