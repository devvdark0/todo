package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devvdark0/todo/internal/model"
	"github.com/devvdark0/todo/internal/service"
)

type TodoHandler struct {
	todoService service.TodoService
}

func NewHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: *service}
}

func (h *TodoHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todoService.ListTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(&tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) GetTask(w http.ResponseWriter, r *http.Request) {

}

func (h *TodoHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := h.todoService.CreateTask(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {

}

func (h *TodoHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

}
