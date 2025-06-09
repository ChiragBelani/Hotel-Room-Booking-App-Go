package booking

import (
	"time"
)

// Booking represents a hotel room booking
type Booking struct {
	Username  string
	Name      string
	Email     string
	RoomType  string // Single / Double / Suite
	RoomNo    int
	// CheckIn   string // Format: YYYY-MM-DD
	// CheckOut  string // Format: YYYY-MM-DD


	CheckIn  time.Time
	CheckOut time.Time
	Nights	int // Number of nights booked
}
