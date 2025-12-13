package auth

import (
	"encoding/json"
	"errors"
	"github.com/devvdark0/todo/internal/model"
	"github.com/devvdark0/todo/internal/service/auth"
	"go.uber.org/zap"
	"net/http"
)

type JWTHandler struct {
	JWTService auth.JWTService
	log        *zap.Logger
}

func NewJWTHandler(service auth.JWTService, log *zap.Logger) *JWTHandler {
	return &JWTHandler{
		JWTService: service,
		log:        log,
	}
}

func (j *JWTHandler) Register(w http.ResponseWriter, r *http.Request) {
	j.log.Info("start proceeding register user request", zap.String("path", r.URL.Path))

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.log.Error("failed to decode request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := j.JWTService.Register(req.Email, req.Username, req.Password); err != nil {
		if errors.Is(err, auth.ErrEmailInUse) {
			j.log.Error("email in use error", zap.Error(err))
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		j.log.Error("failed to create user error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
