package validate

import (
	"fmt"
	"time"
)

func CalculateNights(checkIn, checkOut time.Time) int{
	
	// Parse the dates using time in hours.Parse with the correct layout
	checkIn = checkIn.Truncate(24 * time.Hour)
	checkOut = checkOut.Truncate(24 * time.Hour)

	diff := checkOut.Sub(checkIn)
	days := int(diff.Hours() / 24)

	fmt.Println("Difference in days:", days)
	return days
}