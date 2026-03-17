// Package repository provides data access layer implementations.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thetoppal/encore/internal/model"
)

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// CreateEventParams holds parameters for creating a new event.
type CreateEventParams struct {
	Title          string
	Description    *string
	CategoryID     uuid.UUID
	OrganizerID    uuid.UUID
	ImageURL       *string
	AgeRestriction *int16
}

// UpdateEventParams holds parameters for updating an existing event.
type UpdateEventParams struct {
	Title          string
	Description    *string
	CategoryID     uuid.UUID
	OrganizerID    uuid.UUID
	ImageURL       *string
	AgeRestriction *int16
	Status         model.EventStatus
}

// EventRepository defines the contract for event data access.
type EventRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Event, error)
	List(ctx context.Context, limit, offset int) ([]model.Event, error)
	Create(ctx context.Context, params CreateEventParams) (model.Event, error)
	Update(ctx context.Context, id uuid.UUID, params UpdateEventParams) (model.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// eventRepo is a PostgreSQL implementation of EventRepository.
type eventRepo struct {
	pool *pgxpool.Pool
}

// NewEventRepository creates a new PostgreSQL-backed EventRepository.
func NewEventRepository(pool *pgxpool.Pool) EventRepository {
	return &eventRepo{pool: pool}
}

// GetByID retrieves a single event by its UUID.
// Returns ErrNotFound if no event with the given ID exists.
func (r *eventRepo) GetByID(ctx context.Context, id uuid.UUID) (model.Event, error) {
	query := `
		SELECT id, title, description, category_id, organizer_id,
		       image_url, age_restriction, status, created_at, updated_at
		FROM events
		WHERE id = $1`

	var e model.Event

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.Title, &e.Description, &e.CategoryID, &e.OrganizerID,
		&e.ImageURL, &e.AgeRestriction, &e.Status, &e.CreatedAt, &e.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Event{}, ErrNotFound
	}
	if err != nil {
		return model.Event{}, fmt.Errorf("get event by id: %w", err)
	}

	return e, nil
}

// List returns a paginated list of events ordered by creation date (newest first).
func (r *eventRepo) List(ctx context.Context, limit, offset int) ([]model.Event, error) {
	query := `
		SELECT id, title, description, category_id, organizer_id,
		       image_url, age_restriction, status, created_at, updated_at
		FROM events
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()

	var events []model.Event

	for rows.Next() {
		var e model.Event

		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.CategoryID, &e.OrganizerID,
			&e.ImageURL, &e.AgeRestriction, &e.Status, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan event row: %w", err)
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate event rows: %w", err)
	}

	return events, nil
}

// Create inserts a new event and returns it with server-generated fields.
func (r *eventRepo) Create(ctx context.Context, params CreateEventParams) (model.Event, error) {
	query := `
		INSERT INTO events (title, description, category_id, organizer_id, image_url, age_restriction)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, category_id, organizer_id,
		          image_url, age_restriction, status, created_at, updated_at`

	var e model.Event

	err := r.pool.QueryRow(ctx, query,
		params.Title, params.Description, params.CategoryID, params.OrganizerID,
		params.ImageURL, params.AgeRestriction,
	).Scan(
		&e.ID, &e.Title, &e.Description, &e.CategoryID, &e.OrganizerID,
		&e.ImageURL, &e.AgeRestriction, &e.Status, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return model.Event{}, fmt.Errorf("create event: %w", err)
	}

	return e, nil
}

// Update modifies an existing event and returns the updated version.
// Returns ErrNotFound if no event with the given ID exists.
func (r *eventRepo) Update(ctx context.Context, id uuid.UUID, params UpdateEventParams) (model.Event, error) {
	query := `
		UPDATE events
		SET title = $2, description = $3, category_id = $4, organizer_id = $5,
		    image_url = $6, age_restriction = $7, status = $8, updated_at = now()
		WHERE id = $1
		RETURNING id, title, description, category_id, organizer_id,
		          image_url, age_restriction, status, created_at, updated_at`

	var e model.Event

	err := r.pool.QueryRow(ctx, query, id,
		params.Title, params.Description, params.CategoryID, params.OrganizerID,
		params.ImageURL, params.AgeRestriction, params.Status,
	).Scan(
		&e.ID, &e.Title, &e.Description, &e.CategoryID, &e.OrganizerID,
		&e.ImageURL, &e.AgeRestriction, &e.Status, &e.CreatedAt, &e.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Event{}, ErrNotFound
	}
	if err != nil {
		return model.Event{}, fmt.Errorf("update event: %w", err)
	}

	return e, nil
}

// Delete removes an event by ID. Returns ErrNotFound if the event does not exist.
func (r *eventRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
