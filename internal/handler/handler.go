package handler

import (
	"net/http"

	"github.com/devvdark0/todo/internal/service"
)

type TodoHandler struct {
	todoService service.TodoService
}

func NewHandler(service service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: service}
}

func (h *TodoHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todoService.ListTasks()
	if err != nil {

	}
}

func (h *TodoHandler) GetTask(w http.ResponseWriter, r *http.Request) {

}

func (h *TodoHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

}

func (h *TodoHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {

}

func (h *TodoHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

}
