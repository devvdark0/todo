package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/devvdark0/todo/internal/model"
	"github.com/devvdark0/todo/internal/service"
	"github.com/gorilla/mux"
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

	userID := r.Context().Value("userId").(string)

	tasks, err := h.todoService.ListTasks(userID)
	if err != nil {
		h.log.Error("failed to get tasks", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
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
	taskID := mux.Vars(r)["task_id"]
	userID := r.Context().Value("userId").(string)

	task, err := h.todoService.GetTaskByID(taskID, userID)
	if err != nil {
		h.log.Error("failed to get task with such id", zap.Error(err), zap.String("id", taskID))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
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

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		h.log.Error("unauthorized")
		http.Error(w, "you are not login into service", http.StatusUnauthorized)
		return
	}

	req.UserID = userId

	if err := h.todoService.CreateTask(req); err != nil {
		h.log.Error("failed to create user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *TodoHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info(
		"start update task request",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)
	taskID := mux.Vars(r)["task_id"]
	userID, ok := r.Context().Value("userId").(string)
	if !ok {
		h.log.Error("unauthorized")
		http.Error(w, "you are not log in into service", http.StatusUnauthorized)
		return
	}

	var req model.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode requested body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.todoService.UpdateTask(taskID, userID, req); err != nil {
		h.log.Error("failed to update task", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info(
		"start delete task request",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)
	taskID := mux.Vars(r)["task_id"]
	userID, ok := r.Context().Value("userId").(string)
	if !ok {
		h.log.Error("unauthorized")
		http.Error(w, "you are not logged in", http.StatusUnauthorized)
		return
	}

	if err := h.todoService.DeleteTask(taskID, userID); err != nil {
		h.log.Error("failed to delete task", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
