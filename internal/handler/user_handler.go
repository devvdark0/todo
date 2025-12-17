package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devvdark0/todo/internal/middleware"
	"github.com/devvdark0/todo/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	userStore service.UserStorage
	log       *zap.Logger
}

func NewUserHandler(store service.UserStorage) *UserHandler {
	return &UserHandler{userStore: store}
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := u.userStore.GetByID(userId)
	if err != nil {
		u.log.Error("failed to get user", zap.Error(err))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		u.log.Error("failed to encode data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
