package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"

	"github.com/devvdark0/todo/internal/model"
	"github.com/devvdark0/todo/internal/service"
)

type TodoHandler struct {
	todoService service.TodoService
	log         *zap.Logger
}

func NewHandler(service *service.TodoService, log *zap.Logger) *TodoHandler {
	return &TodoHandler{todoService: *service, log: log}
}

func (h *TodoHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	h.log.Info(
		"start get tasks request",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)
	tasks, err := h.todoService.ListTasks()
	if err != nil {
		h.log.Error("failed to get tasks", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(&tasks); err != nil {
		h.log.Error("failed to encode tasks into json", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info(
		"start get task request",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)
	id := chi.URLParam(r, "id")

	tasks, err := h.todoService.GetTaskByID(id)
	if err != nil {
		h.log.Error("failed to get task with such id", zap.Error(err), zap.String("id", id))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		h.log.Error("failed to encode task into json", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info(
		"request for creating task is started",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)
	var req model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := h.todoService.CreateTask(req); err != nil {
		h.log.Error("failed to create user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info(
		"start update task request",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)
	id := chi.URLParam(r, "id")

	var req model.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode requested body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := h.todoService.UpdateTask(id, req); err != nil {
		h.log.Error("failed to update task", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (h *TodoHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info(
		"start delete task request",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)
	id := chi.URLParam(r, "id")

	if err := h.todoService.DeleteTask(id); err != nil {
		h.log.Error("failed to delete task", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
