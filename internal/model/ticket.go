package model

import (
	"time"

	"github.com/google/uuid"
)

// TicketStatus represents the lifecycle of a ticket.
type TicketStatus string

const (
	TicketStatusActive    TicketStatus = "active"
	TicketStatusUsed      TicketStatus = "used"
	TicketStatusCancelled TicketStatus = "cancelled"
)

// Ticket is a single entry pass. Price is a snapshot at purchase time -
// if the price category changes later, existing tickets keep their price.
type Ticket struct {
	ID              uuid.UUID    `json:"id"`
	OrderID         uuid.UUID    `json:"order_id"`
	SessionID       uuid.UUID    `json:"session_id"`
	PriceCategoryID uuid.UUID    `json:"price_category_id"`
	SeatID          *uuid.UUID   `json:"seat_id,omitempty"` // nil for general_admission
	Price           float64      `json:"price"`
	Status          TicketStatus `json:"status"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}
