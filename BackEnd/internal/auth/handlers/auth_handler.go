package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	authMiddleware "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/middleware"
	authModels "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/models"
	authQueries "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/queries"
	authService "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/service"
)

type AuthHandler struct {
	Service *authService.AuthService
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authModels.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.Service.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidAuthPayload):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, authQueries.ErrEmailAlreadyUsed):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authModels.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.Service.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidLoginPayload):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, authService.ErrInvalidCredentials):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := authMiddleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, authQueries.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
