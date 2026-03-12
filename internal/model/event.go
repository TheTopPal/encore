package model

import (
	"time"

	"github.com/google/uuid"
)

// EventStatus represents the lifecycle state of an event.
type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
)

// Event is a concert, show, or any other bookable occasion.
// One event can have multiple sessions (different dates/venues).
type Event struct {
	ID             uuid.UUID   `json:"id"`
	Title          string      `json:"title"`
	Description    *string     `json:"description,omitempty"`
	CategoryID     uuid.UUID   `json:"category_id"`
	OrganizerID    uuid.UUID   `json:"organizer_id"`
	ImageURL       *string     `json:"image_url,omitempty"`
	AgeRestriction *int16      `json:"age_restriction,omitempty"`
	Status         EventStatus `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
