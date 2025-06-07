package booking

// Booking represents a hotel room booking
type Booking struct {
	Username  string
	Name      string
	Email     string
	RoomType  string // Single / Double / Suite
	RoomNo    int
	Nights    int
}
