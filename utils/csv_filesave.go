package validate

import (
	booking "Hotel-Booking-System/models"
	"fmt"
	"os"
	"encoding/csv"
)

func SaveBookingsToCSV(bookings map[int]booking.Booking, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Username", "Name", "Email", "RoomType", "RoomNo", "CheckIn", "CheckOut", "Nights"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header to CSV: %v", err)
	}

	// Write bookings
	for _, booking := range bookings {
		record := []string{
			booking.Username,
			booking.Name,
			booking.Email,
			booking.RoomType,
			fmt.Sprintf("%d", booking.RoomNo),
			booking.CheckIn.Format("2006-01-02"),
			booking.CheckOut.Format("2006-01-02"),
			fmt.Sprintf("%d", booking.Nights),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record to CSV: %v", err)
		}
	}

	return nil
}