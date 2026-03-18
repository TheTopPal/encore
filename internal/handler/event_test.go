package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/thetoppal/encore/internal/model"
	"github.com/thetoppal/encore/internal/repository"
	"github.com/thetoppal/encore/internal/service"
)

// mockEventRepo is a minimal test double used to drive the real service layer
// from handler tests. Only the methods needed for each test are set.
type mockEventRepo struct {
	getByIDFn func(ctx context.Context, id uuid.UUID) (model.Event, error)
	listFn    func(ctx context.Context, limit, offset int) ([]model.Event, error)
	createFn  func(ctx context.Context, params repository.CreateEventParams) (model.Event, error)
	updateFn  func(ctx context.Context, id uuid.UUID, params repository.UpdateEventParams) (model.Event, error)
	deleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockEventRepo) GetByID(ctx context.Context, id uuid.UUID) (model.Event, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockEventRepo) List(ctx context.Context, limit, offset int) ([]model.Event, error) {
	return m.listFn(ctx, limit, offset)
}

func (m *mockEventRepo) Create(ctx context.Context, params repository.CreateEventParams) (model.Event, error) {
	return m.createFn(ctx, params)
}

func (m *mockEventRepo) Update(ctx context.Context, id uuid.UUID, params repository.UpdateEventParams) (model.Event, error) {
	return m.updateFn(ctx, id, params)
}

func (m *mockEventRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}

// newTestHandler creates an EventHandler backed by the given mock repo.
func newTestHandler(repo repository.EventRepository) *EventHandler {
	svc := service.NewEventService(repo)
	return NewEventHandler(svc)
}

// newChiRequest creates an http.Request with chi URL params injected.
func newChiRequest(method, path string, body []byte, params map[string]string) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if len(params) > 0 {
		rctx := chi.NewRouteContext()
		for k, v := range params {
			rctx.URLParams.Add(k, v)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}

	return req
}

func TestEventHandler_GetByID(t *testing.T) {
	eventID := uuid.New()

	tests := []struct {
		name       string
		id         string
		repoFn     func(context.Context, uuid.UUID) (model.Event, error)
		wantStatus int
	}{
		{
			name:       "invalid uuid",
			id:         "not-a-uuid",
			repoFn:     nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			id:   eventID.String(),
			repoFn: func(_ context.Context, _ uuid.UUID) (model.Event, error) {
				return model.Event{}, repository.ErrNotFound
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "success",
			id:   eventID.String(),
			repoFn: func(_ context.Context, id uuid.UUID) (model.Event, error) {
				return model.Event{ID: id, Title: "Concert"}, nil
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockEventRepo{getByIDFn: tt.repoFn}
			h := newTestHandler(repo)

			req := newChiRequest(http.MethodGet, "/events/"+tt.id, nil, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()

			h.GetByID(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestEventHandler_List(t *testing.T) {
	repo := &mockEventRepo{
		listFn: func(_ context.Context, _, _ int) ([]model.Event, error) {
			return []model.Event{{Title: "A"}, {Title: "B"}}, nil
		},
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/events?limit=10&offset=0", nil)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rr.Code, http.StatusOK)
	}

	var events []model.Event
	if err := json.NewDecoder(rr.Body).Decode(&events); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("got %d events, want 2", len(events))
	}
}

func TestEventHandler_Create(t *testing.T) {
	catID := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name:       "invalid json",
			body:       "not json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty title - validation error",
			body:       createEventRequest{Title: "", CategoryID: catID, OrganizerID: orgID},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "valid request",
			body:       createEventRequest{Title: "Concert", CategoryID: catID, OrganizerID: orgID},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockEventRepo{
				createFn: func(_ context.Context, params repository.CreateEventParams) (model.Event, error) {
					return model.Event{
						ID:          uuid.New(),
						Title:       params.Title,
						CategoryID:  params.CategoryID,
						OrganizerID: params.OrganizerID,
						Status:      model.EventStatusDraft,
					}, nil
				},
			}
			h := newTestHandler(repo)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			req := newChiRequest(http.MethodPost, "/events", bodyBytes, nil)
			rr := httptest.NewRecorder()

			h.Create(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d; body: %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
		})
	}
}

func TestEventHandler_Delete(t *testing.T) {
	eventID := uuid.New()

	tests := []struct {
		name       string
		id         string
		repoErr    error
		wantStatus int
	}{
		{
			name:       "invalid uuid",
			id:         "bad",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			id:         eventID.String(),
			repoErr:    repository.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "success",
			id:         eventID.String(),
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockEventRepo{
				deleteFn: func(_ context.Context, _ uuid.UUID) error {
					return tt.repoErr
				},
			}
			h := newTestHandler(repo)

			req := newChiRequest(http.MethodDelete, "/events/"+tt.id, nil, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()

			h.Delete(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}
