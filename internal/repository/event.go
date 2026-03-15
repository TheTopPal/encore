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

// EventRepository defines the contract for event data access.
type EventRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Event, error)
	List(ctx context.Context, limit, offset int) ([]model.Event, error)
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
