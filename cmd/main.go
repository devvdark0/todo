package main

import (
	"github.com/devvdark0/todo/internal/config"
	"log"
	"net/http"

	"github.com/devvdark0/todo/internal/handler"
	"github.com/devvdark0/todo/internal/service"
	"github.com/devvdark0/todo/internal/storage"
	"github.com/devvdark0/todo/pkg/db"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	log.Println("successfully connected to mariadb!")

	//TODO: init logger

	store := storage.NewStore(database)
	todoService := service.NewService(store)
	todoHandler := handler.NewHandler(todoService)

	router := configureRouter(todoHandler)

	srv := http.Server{
		Addr:         "localhost:" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Println("starting server on port :80...")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
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
