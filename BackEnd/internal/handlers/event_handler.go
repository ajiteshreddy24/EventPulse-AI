package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	authMiddleware "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/middleware"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/service"
	"github.com/gorilla/mux"
)

type EventHandler struct {
	Service *service.EventService
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
	events, err := h.Service.GetEvents()
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

	event, err := h.Service.GetEventByID(id)
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
