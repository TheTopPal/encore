package model

import (
	"time"

	"github.com/google/uuid"
)

// EventSessionStatus represents the lifecycle of a session.
type EventSessionStatus string

const (
	EventSessionStatusScheduled EventSessionStatus = "scheduled"
	EventSessionStatusSelling   EventSessionStatus = "selling"
	EventSessionStatusSoldOut   EventSessionStatus = "sold_out"
	EventSessionStatusCancelled EventSessionStatus = "cancelled"
)

// EventSession is a specific showing of an event (date + venue hall).
type EventSession struct {
	ID        uuid.UUID          `json:"id"`
	EventID   uuid.UUID          `json:"event_id"`
	HallID    uuid.UUID          `json:"hall_id"`
	StartsAt  time.Time          `json:"starts_at"`
	EndsAt    time.Time          `json:"ends_at"`
	Status    EventSessionStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// EventCategory is a reference table (concert, theatre, standup...).
type EventCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// PriceCategory defines pricing tiers for a session (VIP, Standard...).
type PriceCategory struct {
	ID            uuid.UUID `json:"id"`
	SessionID     uuid.UUID `json:"session_id"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	TotalQuantity *int      `json:"total_quantity,omitempty"` // only for general_admission
	Section       *string   `json:"section,omitempty"`
}
