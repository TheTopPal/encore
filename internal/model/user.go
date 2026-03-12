package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRole determines what actions a user can perform in the system.
type UserRole string

const (
	UserRoleBuyer     UserRole = "buyer"
	UserRoleAdmin     UserRole = "admin"
	UserRoleOrganizer UserRole = "organizer"
)

// User represents any person in the system: buyer, organizer, or admin.
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never expose in JSON
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Username     string    `json:"username"`
	Role         UserRole  `json:"role"`
	Phone        *string   `json:"phone,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
