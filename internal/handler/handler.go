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

}
