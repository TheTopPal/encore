package model

import (
	"time"

	"github.com/google/uuid"
)

// Venue is a physical location (concert hall, stadium, club).
// Owned by an organizer.
type Venue struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	OrganizerID uuid.UUID `json:"organizer_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SeatingType defines how seats are allocated in a hall.
type SeatingType string

const (
	SeatingTypeAssigned         SeatingType = "assigned"
	SeatingTypeGeneralAdmission SeatingType = "general_admission"
)

// VenueHall is a hall inside a venue (e.g. "Main Hall", "Hall 1").
type VenueHall struct {
	ID          uuid.UUID   `json:"id"`
	VenueID     uuid.UUID   `json:"venue_id"`
	Name        string      `json:"name"`
	SeatingType SeatingType `json:"seating_type"`
	Capacity    int         `json:"capacity"`
}

// Seat is a specific seat in a hall. Only exists for assigned seating.
// For general_admission halls, tickets have no seat reference.
type Seat struct {
	ID         uuid.UUID `json:"id"`
	HallID     uuid.UUID `json:"hall_id"`
	RowNumber  string    `json:"row_number"`
	SeatNumber int       `json:"seat_number"`
	Section    *string   `json:"section,omitempty"`
}
