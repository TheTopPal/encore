package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/thetoppal/encore/internal/model"
	"github.com/thetoppal/encore/internal/repository"
)

// mockEventRepo is a test double for repository.EventRepository.
// Each method delegates to a function field, so individual tests
// can configure only the behavior they care about.
type mockEventRepo struct {
	getByIDFn func(ctx context.Context, id uuid.UUID) (model.Event, error)
	listFn    func(ctx context.Context, filter repository.EventFilter, limit, offset int) ([]model.Event, int, error)
	createFn  func(ctx context.Context, params repository.CreateEventParams) (model.Event, error)
	updateFn  func(ctx context.Context, id uuid.UUID, params repository.UpdateEventParams) (model.Event, error)
	deleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockEventRepo) GetByID(ctx context.Context, id uuid.UUID) (model.Event, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockEventRepo) List(ctx context.Context, filter repository.EventFilter, limit, offset int) ([]model.Event, int, error) {
	return m.listFn(ctx, filter, limit, offset)
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
