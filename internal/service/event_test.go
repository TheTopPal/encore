package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/thetoppal/encore/internal/model"
	"github.com/thetoppal/encore/internal/repository"
)

func TestEventService_Create(t *testing.T) {
	catID := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name    string
		input   CreateEventInput
		wantErr error
	}{
		{
			name:    "empty title",
			input:   CreateEventInput{Title: "", CategoryID: catID, OrganizerID: orgID},
			wantErr: ErrValidation,
		},
		{
			name:    "nil category_id",
			input:   CreateEventInput{Title: "Concert", CategoryID: uuid.Nil, OrganizerID: orgID},
			wantErr: ErrValidation,
		},
		{
			name:    "nil organizer_id",
			input:   CreateEventInput{Title: "Concert", CategoryID: catID, OrganizerID: uuid.Nil},
			wantErr: ErrValidation,
		},
		{
			name:  "valid input",
			input: CreateEventInput{Title: "Concert", CategoryID: catID, OrganizerID: orgID},
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
			svc := NewEventService(repo)

			event, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got err %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if event.Title != tt.input.Title {
				t.Errorf("got title %q, want %q", event.Title, tt.input.Title)
			}
		})
	}
}

func TestEventService_Update(t *testing.T) {
	eventID := uuid.New()
	catID := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name    string
		input   UpdateEventInput
		repoErr error
		wantErr error
	}{
		{
			name:    "empty title",
			input:   UpdateEventInput{Title: "", CategoryID: catID, OrganizerID: orgID, Status: model.EventStatusDraft},
			wantErr: ErrValidation,
		},
		{
			name:    "invalid status",
			input:   UpdateEventInput{Title: "Concert", CategoryID: catID, OrganizerID: orgID, Status: "bogus"},
			wantErr: ErrValidation,
		},
		{
			name:    "not found",
			input:   UpdateEventInput{Title: "Concert", CategoryID: catID, OrganizerID: orgID, Status: model.EventStatusPublished},
			repoErr: repository.ErrNotFound,
			wantErr: repository.ErrNotFound,
		},
		{
			name:  "valid update",
			input: UpdateEventInput{Title: "Updated", CategoryID: catID, OrganizerID: orgID, Status: model.EventStatusPublished},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockEventRepo{
				updateFn: func(_ context.Context, _ uuid.UUID, params repository.UpdateEventParams) (model.Event, error) {
					if tt.repoErr != nil {
						return model.Event{}, tt.repoErr
					}
					return model.Event{
						ID:     eventID,
						Title:  params.Title,
						Status: params.Status,
					}, nil
				},
			}
			svc := NewEventService(repo)

			event, err := svc.Update(context.Background(), eventID, tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got err %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if event.Title != tt.input.Title {
				t.Errorf("got title %q, want %q", event.Title, tt.input.Title)
			}
		})
	}
}

func TestEventService_List(t *testing.T) {
	published := model.EventStatusPublished
	catID := uuid.New()

	tests := []struct {
		name       string
		filter     EventListFilter
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
		wantTotal  int
		wantErr    error
	}{
		{name: "defaults for zero", limit: 0, offset: 0, wantLimit: 20, wantOffset: 0, wantTotal: 5},
		{name: "negative limit", limit: -5, offset: 0, wantLimit: 20, wantOffset: 0, wantTotal: 5},
		{name: "exceeds max", limit: 200, offset: 0, wantLimit: 100, wantOffset: 0, wantTotal: 5},
		{name: "negative offset", limit: 10, offset: -3, wantLimit: 10, wantOffset: 0, wantTotal: 5},
		{name: "normal values", limit: 50, offset: 10, wantLimit: 50, wantOffset: 10, wantTotal: 5},
		{
			name:      "filter by status",
			filter:    EventListFilter{Status: &published},
			limit:     10,
			offset:    0,
			wantLimit: 10, wantOffset: 0, wantTotal: 5,
		},
		{
			name:      "filter by category_id",
			filter:    EventListFilter{CategoryID: &catID},
			limit:     10,
			offset:    0,
			wantLimit: 10, wantOffset: 0, wantTotal: 5,
		},
		{
			name:    "invalid status filter",
			filter:  EventListFilter{Status: ptr(model.EventStatus("bogus"))},
			limit:   10,
			offset:  0,
			wantErr: ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotLimit, gotOffset int

			repo := &mockEventRepo{
				listFn: func(_ context.Context, _ repository.EventFilter, limit, offset int) ([]model.Event, int, error) {
					gotLimit = limit
					gotOffset = offset
					return []model.Event{}, 5, nil
				},
			}
			svc := NewEventService(repo)

			events, total, err := svc.List(context.Background(), tt.filter, tt.limit, tt.offset)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got err %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("got limit %d, want %d", gotLimit, tt.wantLimit)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("got offset %d, want %d", gotOffset, tt.wantOffset)
			}
			if total != tt.wantTotal {
				t.Errorf("got total %d, want %d", total, tt.wantTotal)
			}
			_ = events
		})
	}
}

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}

func TestEventService_GetByID(t *testing.T) {
	eventID := uuid.New()

	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{name: "found", repoErr: nil},
		{name: "not found", repoErr: repository.ErrNotFound, wantErr: repository.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockEventRepo{
				getByIDFn: func(_ context.Context, id uuid.UUID) (model.Event, error) {
					if tt.repoErr != nil {
						return model.Event{}, tt.repoErr
					}
					return model.Event{ID: id, Title: "Concert"}, nil
				},
			}
			svc := NewEventService(repo)

			event, err := svc.GetByID(context.Background(), eventID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got err %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if event.ID != eventID {
				t.Errorf("got id %v, want %v", event.ID, eventID)
			}
		})
	}
}

func TestEventService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{name: "success", repoErr: nil},
		{name: "not found", repoErr: repository.ErrNotFound, wantErr: repository.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockEventRepo{
				deleteFn: func(_ context.Context, _ uuid.UUID) error {
					return tt.repoErr
				},
			}
			svc := NewEventService(repo)

			err := svc.Delete(context.Background(), uuid.New())

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got err %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
