package validate

import (
	"fmt"
	"strings"
)

// ValidateBooking checks if the booking details are valid
func InputValidation(name, email, roomType string, nights int) bool {
	if len(name) < 2 {
		fmt.Println("Name must be more than 2 characters.")
		return false
	}
	if !(len(email) > 5 && strings.Contains(email, "@") && (email[len(email)-4:] == ".com" || email[len(email)-4:] == ".org" || email[len(email)-4:] == ".net")){
		fmt.Println("Email must be valid and contain '@' and a domain like '.com', '.org', or '.net'.")
		return false
	}
	if roomType != "Single" && roomType != "Double" && roomType != "Suite" {
		fmt.Println("Room type must be either Single, Double, or Suite.")
		return false
	}
	if nights <= 0 {
		fmt.Println("Number of nights must be greater than 0.")
		return false
	}
	return true
}
