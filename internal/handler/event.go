// Package handler provides HTTP handlers for the API.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/thetoppal/encore/internal/repository"
	"github.com/thetoppal/encore/internal/service"
)

// EventHandler handles HTTP requests for events.
type EventHandler struct {
	service *service.EventService
}

// NewEventHandler creates a new EventHandler.
func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{service: svc}
}

// Routes registers event endpoints on the given router.
func (h *EventHandler) Routes(r chi.Router) {
	r.Get("/events", h.List)
	r.Get("/events/{id}", h.GetByID)
}

// GetByID handles GET /events/{id}.
func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid event id"})
		return
	}

	event, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, errorResponse{"event not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{"internal error"})
		return
	}

	writeJSON(w, http.StatusOK, event)
}

// List handles GET /events?limit=20&offset=0.
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 20)
	offset := queryInt(r, "offset", 0)

	events, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{"internal error"})
		return
	}

	writeJSON(w, http.StatusOK, events)
}

// errorResponse is a standard error payload.
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON encodes v as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck // best-effort response write
}

// queryInt reads an integer query parameter with a default fallback.
func queryInt(r *http.Request, key string, defaultVal int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}

	return val
}
