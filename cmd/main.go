package main

import (
	"log"
	"net/http"

	"github.com/devvdark0/todo/internal/handler"
	"github.com/devvdark0/todo/internal/storage"
	"github.com/devvdark0/todo/pkg/db"
	"github.com/go-chi/chi/v5"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	log.Println("successfully connected to mariadb!")

	store := storage.NewStore(database)
	_ = store

	log.Println("starting server on port :80...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func configureRouter(todoHandler *handler.TodoHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/tasks", todoHandler.GetTasks)
}
