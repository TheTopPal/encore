// Package handler provides HTTP handlers for the API.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/thetoppal/encore/internal/model"
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
	r.Post("/events", h.Create)
	r.Put("/events/{id}", h.Update)
	r.Delete("/events/{id}", h.Delete)
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

// eventListResponse wraps a page of events with the total matching count.
type eventListResponse struct {
	Items []model.Event `json:"items"`
	Total int           `json:"total"`
}

// List handles GET /events?limit=20&offset=0&status=published&category_id=...
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 20)
	offset := queryInt(r, "offset", 0)

	var filter service.EventListFilter

	if s := r.URL.Query().Get("status"); s != "" {
		status := model.EventStatus(s)
		filter.Status = &status
	}

	if s := r.URL.Query().Get("category_id"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{"invalid category_id"})
			return
		}
		filter.CategoryID = &id
	}

	events, total, err := h.service.List(r.Context(), filter, limit, offset)
	if errors.Is(err, service.ErrValidation) {
		writeJSON(w, http.StatusBadRequest, errorResponse{err.Error()})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{"internal error"})
		return
	}

	if events == nil {
		events = []model.Event{}
	}

	writeJSON(w, http.StatusOK, eventListResponse{Items: events, Total: total})
}

// createEventRequest is the JSON body for POST /events.
type createEventRequest struct {
	Title          string    `json:"title"`
	Description    *string   `json:"description,omitempty"`
	CategoryID     uuid.UUID `json:"category_id"`
	OrganizerID    uuid.UUID `json:"organizer_id"`
	ImageURL       *string   `json:"image_url,omitempty"`
	AgeRestriction *int16    `json:"age_restriction,omitempty"`
}

// updateEventRequest is the JSON body for PUT /events/{id}.
type updateEventRequest struct {
	Title          string            `json:"title"`
	Description    *string           `json:"description,omitempty"`
	CategoryID     uuid.UUID         `json:"category_id"`
	OrganizerID    uuid.UUID         `json:"organizer_id"`
	ImageURL       *string           `json:"image_url,omitempty"`
	AgeRestriction *int16            `json:"age_restriction,omitempty"`
	Status         model.EventStatus `json:"status"`
}

// Create handles POST /events.
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid request body"})
		return
	}

	event, err := h.service.Create(r.Context(), service.CreateEventInput{
		Title:          req.Title,
		Description:    req.Description,
		CategoryID:     req.CategoryID,
		OrganizerID:    req.OrganizerID,
		ImageURL:       req.ImageURL,
		AgeRestriction: req.AgeRestriction,
	})
	if errors.Is(err, service.ErrValidation) {
		writeJSON(w, http.StatusBadRequest, errorResponse{err.Error()})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{"internal error"})
		return
	}

	writeJSON(w, http.StatusCreated, event)
}

// Update handles PUT /events/{id}.
func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid event id"})
		return
	}

	var req updateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid request body"})
		return
	}

	event, err := h.service.Update(r.Context(), id, service.UpdateEventInput{
		Title:          req.Title,
		Description:    req.Description,
		CategoryID:     req.CategoryID,
		OrganizerID:    req.OrganizerID,
		ImageURL:       req.ImageURL,
		AgeRestriction: req.AgeRestriction,
		Status:         req.Status,
	})
	if errors.Is(err, service.ErrValidation) {
		writeJSON(w, http.StatusBadRequest, errorResponse{err.Error()})
		return
	}
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

// Delete handles DELETE /events/{id}.
func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid event id"})
		return
	}

	err = h.service.Delete(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, errorResponse{"event not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{"internal error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
