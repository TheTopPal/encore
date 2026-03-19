// Package repository provides data access layer implementations.
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

// EventFilter holds optional filter criteria for listing events.
type EventFilter struct {
	Status     *model.EventStatus
	CategoryID *uuid.UUID
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
	List(ctx context.Context, filter EventFilter, limit, offset int) ([]model.Event, int, error)
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

// List returns a filtered, paginated list of events ordered by creation date
// (newest first) together with the total number of matching events.
func (r *eventRepo) List(ctx context.Context, filter EventFilter, limit, offset int) ([]model.Event, int, error) {
	where, args := buildEventFilter(filter)
	argIdx := len(args) + 1

	// Total count of matching events (ignoring LIMIT/OFFSET).
	countQuery := "SELECT COUNT(*) FROM events " + where

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count events: %w", err)
	}

	// Data page.
	dataQuery := fmt.Sprintf(`
		SELECT id, title, description, category_id, organizer_id,
		       image_url, age_restriction, status, created_at, updated_at
		FROM events
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)

	dataArgs := make([]any, 0, len(args)+2)
	dataArgs = append(dataArgs, args...)
	dataArgs = append(dataArgs, limit, offset)

	rows, err := r.pool.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()

	var events []model.Event

	for rows.Next() {
		var e model.Event

		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.CategoryID, &e.OrganizerID,
			&e.ImageURL, &e.AgeRestriction, &e.Status, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan event row: %w", err)
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate event rows: %w", err)
	}

	return events, total, nil
}

// buildEventFilter constructs a WHERE clause and parameter list from the filter.
func buildEventFilter(filter EventFilter) (string, []any) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *filter.Status)
		argIdx++
	}

	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argIdx))
		args = append(args, *filter.CategoryID)
	}

	if len(conditions) == 0 {
		return "", nil
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
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
