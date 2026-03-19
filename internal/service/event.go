// Package service implements business logic on top of the repository layer.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/thetoppal/encore/internal/model"
	"github.com/thetoppal/encore/internal/repository"
)

// ErrValidation indicates invalid input data.
var ErrValidation = errors.New("validation error")

// CreateEventInput holds the data needed to create a new event.
type CreateEventInput struct {
	Title          string
	Description    *string
	CategoryID     uuid.UUID
	OrganizerID    uuid.UUID
	ImageURL       *string
	AgeRestriction *int16
}

// EventListFilter holds optional filter criteria for listing events.
type EventListFilter struct {
	Status     *model.EventStatus
	CategoryID *uuid.UUID
}

// UpdateEventInput holds the data needed to update an event.
type UpdateEventInput struct {
	Title          string
	Description    *string
	CategoryID     uuid.UUID
	OrganizerID    uuid.UUID
	ImageURL       *string
	AgeRestriction *int16
	Status         model.EventStatus
}

// EventService provides business operations on events.
type EventService struct {
	repo repository.EventRepository
}

// NewEventService creates a new EventService.
func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

// GetByID returns a single event by ID.
func (s *EventService) GetByID(ctx context.Context, id uuid.UUID) (model.Event, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Event{}, fmt.Errorf("service get event: %w", err)
	}

	return event, nil
}

// List returns a filtered, paginated list of events with total count.
func (s *EventService) List(ctx context.Context, filter EventListFilter, limit, offset int) ([]model.Event, int, error) {
	if filter.Status != nil && !validEventStatus(*filter.Status) {
		return nil, 0, fmt.Errorf("%w: invalid status filter %q", ErrValidation, *filter.Status)
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	repoFilter := repository.EventFilter{
		Status:     filter.Status,
		CategoryID: filter.CategoryID,
	}

	events, total, err := s.repo.List(ctx, repoFilter, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("service list events: %w", err)
	}

	return events, total, nil
}

// Create validates input and creates a new event.
func (s *EventService) Create(ctx context.Context, input CreateEventInput) (model.Event, error) {
	if input.Title == "" {
		return model.Event{}, fmt.Errorf("%w: title is required", ErrValidation)
	}
	if input.CategoryID == uuid.Nil {
		return model.Event{}, fmt.Errorf("%w: category_id is required", ErrValidation)
	}
	if input.OrganizerID == uuid.Nil {
		return model.Event{}, fmt.Errorf("%w: organizer_id is required", ErrValidation)
	}

	event, err := s.repo.Create(ctx, repository.CreateEventParams{
		Title:          input.Title,
		Description:    input.Description,
		CategoryID:     input.CategoryID,
		OrganizerID:    input.OrganizerID,
		ImageURL:       input.ImageURL,
		AgeRestriction: input.AgeRestriction,
	})
	if err != nil {
		return model.Event{}, fmt.Errorf("service create event: %w", err)
	}

	return event, nil
}

// Update validates input and updates an existing event.
func (s *EventService) Update(ctx context.Context, id uuid.UUID, input UpdateEventInput) (model.Event, error) {
	if input.Title == "" {
		return model.Event{}, fmt.Errorf("%w: title is required", ErrValidation)
	}
	if input.CategoryID == uuid.Nil {
		return model.Event{}, fmt.Errorf("%w: category_id is required", ErrValidation)
	}
	if input.OrganizerID == uuid.Nil {
		return model.Event{}, fmt.Errorf("%w: organizer_id is required", ErrValidation)
	}
	if !validEventStatus(input.Status) {
		return model.Event{}, fmt.Errorf("%w: invalid status %q", ErrValidation, input.Status)
	}

	event, err := s.repo.Update(ctx, id, repository.UpdateEventParams{
		Title:          input.Title,
		Description:    input.Description,
		CategoryID:     input.CategoryID,
		OrganizerID:    input.OrganizerID,
		ImageURL:       input.ImageURL,
		AgeRestriction: input.AgeRestriction,
		Status:         input.Status,
	})
	if err != nil {
		return model.Event{}, fmt.Errorf("service update event: %w", err)
	}

	return event, nil
}

// Delete removes an event by ID.
func (s *EventService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service delete event: %w", err)
	}

	return nil
}

// validEventStatus checks if the given status is a known EventStatus value.
func validEventStatus(s model.EventStatus) bool {
	switch s {
	case model.EventStatusDraft, model.EventStatusPublished, model.EventStatusCancelled:
		return true
	default:
		return false
	}
}
