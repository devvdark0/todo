package main

import (
	"log"
	"net/http"

	"github.com/devvdark0/todo/internal/handler"
	"github.com/devvdark0/todo/internal/service"
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
	todoService := service.NewService(store)
	todoHandler := handler.NewHandler(todoService)

	router := configureRouter(todoHandler)

	log.Println("starting server on port :80...")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}

func configureRouter(todoHandler *handler.TodoHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/tasks", todoHandler.GetTasks)
	router.Post("/tasks", todoHandler.CreateTask)
	router.Get("/tasks/{id}", todoHandler.GetTask)
	router.Put("/tasks/{id}", todoHandler.UpdateTask)
	router.Delete("/tasks/{id}", todoHandler.DeleteTask)
	return router
}
