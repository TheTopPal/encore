// Package service implements business logic on top of the repository layer.
package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/thetoppal/encore/internal/model"
	"github.com/thetoppal/encore/internal/repository"
)

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

// List returns a paginated list of events.
func (s *EventService) List(ctx context.Context, limit, offset int) ([]model.Event, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	events, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("service list events: %w", err)
	}

	return events, nil
}
