package app

import (
	"net/http"

	"github.com/devvdark0/todo/internal/config"
	"github.com/devvdark0/todo/internal/handler"
	"github.com/devvdark0/todo/internal/middleware"
	"github.com/devvdark0/todo/internal/service"
	"github.com/devvdark0/todo/internal/storage"
	"github.com/devvdark0/todo/pkg/db"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitApp() error {
	cfg, err := config.MustLoad()
	if err != nil {
		return err
	}

	log := configureLogger(cfg)

	database, err := db.InitDB(cfg)
	if err != nil {
		return err
	}
	defer database.Close()

	taskStore := storage.NewStore(database, log)
	taskService := service.NewService(taskStore)
	taskHandler := handler.NewHandler(taskService, log)

	userStore := storage.NewUserStore(database, log)
	authService := service.NewJWTService([]byte(cfg.JWTConfig.Secret), cfg.JWTConfig.TokenTTL, userStore)
	authHandler := handler.NewJWTHandler(*authService, log)
	userHandler := handler.NewUserHandler(userStore)

	r := configureRouter(taskHandler, authHandler, userHandler, authService)

	srv := http.Server{
		Addr:         "localhost:" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting server on port:", zap.String("port", cfg.Port))
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

	return nil

}

func configureRouter(
	todoHandler *handler.TodoHandler,
	authHandler *handler.JWTHandler,
	userHandler *handler.UserHandler,
	authService *service.JWTService,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	protected := r.PathPrefix("/api").Subrouter()

	protected.Use(middleware.AuthMiddleware(*authService))
	protected.HandleFunc("/tasks", todoHandler.GetTasks).Methods("GET")
	protected.HandleFunc("/tasks/{task_id}", todoHandler.GetTask).Methods("GET")
	protected.HandleFunc("/tasks", todoHandler.CreateTask).Methods("POST")
	protected.HandleFunc("/tasks/{task_id}", todoHandler.UpdateTask).Methods("PUT")
	protected.HandleFunc("/tasks/{task_id}", todoHandler.DeleteTask).Methods("DELETE")
	protected.HandleFunc("/profile", userHandler.Profile).Methods("GET")

	return r
}

func configureLogger(cfg *config.Config) *zap.Logger {
	if cfg.Env == "local" {
		return zap.Must(zap.NewDevelopment())
	}
	return zap.Must(zap.NewProduction())
}
