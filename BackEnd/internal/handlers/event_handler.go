package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	authMiddleware "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/middleware"
	authService "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/service"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/service"
	"github.com/gorilla/mux"
)

type EventHandler struct {
	Service     *service.EventService
	AuthService *authService.AuthService
}

/* =========================
   CREATE EVENT
========================= */

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateEvent(&event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

/* =========================
   GET ALL EVENTS
========================= */

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	userID := h.optionalUserID(r)
	events, err := h.Service.GetEvents(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

/* =========================
   GET EVENT BY ID
========================= */

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	userID := h.optionalUserID(r)
	event, err := h.Service.GetEventByID(id, userID)
	if err != nil {
		if errors.Is(err, queries.ErrEventNotFound) {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(event)
}

/* =========================
   UPDATE EVENT
========================= */

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	event.ID = id

	if err := h.Service.UpdateEvent(&event); err != nil {
		if errors.Is(err, queries.ErrEventNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(event)
}

/* =========================
   DELETE EVENT
========================= */

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteEvent(id); err != nil {
		if errors.Is(err, queries.ErrEventNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* =========================
   RSVP
========================= */

func (h *EventHandler) RSVP(w http.ResponseWriter, r *http.Request) {
	userID, ok := authMiddleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	if err := h.Service.RSVP(userID, eventID); err != nil {
		switch {
		case errors.Is(err, queries.ErrRSVPAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, queries.ErrEventFull):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, queries.ErrEventNotFound):
			http.Error(w, "event not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

/* =========================
   CANCEL RSVP
========================= */

func (h *EventHandler) CancelRSVP(w http.ResponseWriter, r *http.Request) {
	userID, ok := authMiddleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	if err := h.Service.CancelRSVP(userID, eventID); err != nil {
		if errors.Is(err, queries.ErrEventNotFound) {
			http.Error(w, "rsvp not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) JoinWaitlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := authMiddleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	if err := h.Service.JoinWaitlist(userID, eventID); err != nil {
		switch {
		case errors.Is(err, queries.ErrAlreadyOnWaitlist), errors.Is(err, queries.ErrEventFull), errors.Is(err, queries.ErrRSVPAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, queries.ErrEventNotFound):
			http.Error(w, "event not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Added to waitlist"})
}

func (h *EventHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := authMiddleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var payload struct {
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.Service.CreateComment(userID, eventID, payload.Comment)
	if err != nil {
		switch {
		case errors.Is(err, queries.ErrEventNotFound):
			http.Error(w, "event not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *EventHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	comments, err := h.Service.GetComments(eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func (h *EventHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteComment(commentID); err != nil {
		if errors.Is(err, queries.ErrEventNotFound) {
			http.Error(w, "comment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted"})
}

func (h *EventHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID, ok := authMiddleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	events, err := h.Service.GetRecommendations(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) optionalUserID(r *http.Request) *int {
	if h.AuthService == nil {
		return nil
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.AuthService.ParseToken(tokenString)
	if err != nil {
		return nil
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return nil
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		return nil
	}

	return &userID
}
