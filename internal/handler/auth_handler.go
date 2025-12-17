package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/devvdark0/todo/internal/model"
	"github.com/devvdark0/todo/internal/service"
	"go.uber.org/zap"
)

type JWTHandler struct {
	JWTService service.JWTService
	log        *zap.Logger
}

func NewJWTHandler(service service.JWTService, log *zap.Logger) *JWTHandler {
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
		if errors.Is(err, service.ErrEmailInUse) {
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

func (j *JWTHandler) Login(w http.ResponseWriter, r *http.Request) {
	j.log.Info("start proceeding login request", zap.String("path", r.URL.Path))

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.log.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := j.JWTService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			j.log.Error("invalid credentials err", zap.Error(err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		j.log.Error("failed to login", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := model.LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		j.log.Error("failed to encode data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
