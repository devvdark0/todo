package main

import (
	"github.com/devvdark0/todo/internal/config"
	task3 "github.com/devvdark0/todo/internal/handler/task"
	task2 "github.com/devvdark0/todo/internal/service/task"
	"github.com/devvdark0/todo/internal/storage/task"
	"go.uber.org/zap"
	"net/http"

	"github.com/devvdark0/todo/pkg/db"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		panic(err)
	}

	log := configureLogger(cfg)

	database, err := db.InitDB(cfg)
	if err != nil {
		log.Panic("failed to connect to the database: " + err.Error())
	}
	defer database.Close()
	log.Info("successfully connected to mariadb!")

	//TODO: impl log in layers
	store := task.NewStore(database, log)
	todoService := task2.NewService(store)
	todoHandler := task3.NewHandler(todoService, log)

	router := configureRouter(todoHandler)

	srv := http.Server{
		Addr:         "localhost:" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting server on port:", zap.String("port", cfg.Port))
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func configureRouter(todoHandler *task3.TodoHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/tasks", todoHandler.GetTasks)
	router.Post("/tasks", todoHandler.CreateTask)
	router.Get("/tasks/{id}", todoHandler.GetTask)
	router.Put("/tasks/{id}", todoHandler.UpdateTask)
	router.Delete("/tasks/{id}", todoHandler.DeleteTask)
	return router
}

func configureLogger(cfg *config.Config) *zap.Logger {
	if cfg.Env == "local" {
		return zap.Must(zap.NewDevelopment())
	}
	return zap.Must(zap.NewProduction())
}
