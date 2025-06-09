package main

import (
	"fmt"
	// "strings"
	//"os/user"
	"sync"

	// "strings"
	"Hotel-Booking-System/models" // import the models package for Booking struct
	"Hotel-Booking-System/utils"  // name of the folder is to be mentioned here
	"encoding/json"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)



var bookings = loadBookingsFromFile("bookingData.json")

var singleRoomNo = 101
var singleRoomNoIdx = 0
var doubleRoomNo = 201
var doubleRoomNoIdx = 0
var suiteRoomNo = 301
var suiteRoomNoIdx = 0

var singleBookarray = [5]int{0,0,0,0,0}
var doubleBookarray = [5]int{0,0,0,0,0}
var suiteBookarray = [5]int{0,0,0,0,0}


var mu = sync.Mutex{}

func main() {
	var username = userLogin()
	if username == "f"{
		fmt.Println("Login failed. Please try again.")
		os.Exit(0) // Exit the program if login fails
	} else {
	updateBookArrayfromJson()
	var choice int
	for {
		showMenu(username)
		fmt.Print("Please enter your choice: ")
		fmt.Scan(&choice)
		if choice < 1 || choice > 4 {
			fmt.Println("Invalid choice. Please select a valid option (1-4).")
			continue
		}

		switch choice {
		case 1:
			bookRoom(username)
		case 2:
			viewUserBookings(username)
		case 3:
			viewUserBookings(username)
			mu.Lock()
			cancelBooking(username)
		case 4:
			username = userLogin()
		}
	}
	}
}

func showMenu(username string) {
	fmt.Println("\n#########################################################################")
	fmt.Println("🏨 Welcome to the Hotel Booking System")
	fmt.Println("Logged in as:", username)
	fmt.Println("1. Book a Room")
	fmt.Println("2. View Your Bookings")
	fmt.Println("3. Cancel a Booking")
	fmt.Println("4. Logout")
	fmt.Println("#########################################################################")
}

func bookRoom(username string) {
	
	var name, email, roomType string

	var roomNo, nights int
	var checkIn time.Time
	var checkOut time.Time

	fmt.Print("Enter your name: ")
	fmt.Scan(&name)
	fmt.Print("Enter your email: ")
	fmt.Scan(&email)
	fmt.Print("Enter room type (Single/Double/Suite): ")
	fmt.Scan(&roomType)
	
	if !checkAvailableRoom(roomType){
		return
	}

	fmt.Print("Enter the nights: ")
	fmt.Scan(&nights)

	checkIn = time.Now()
	checkOut = checkIn.AddDate(0, 0, nights)
	

	if !validate.InputValidation(name, email, roomType, nights){
		fmt.Println("Invalid input. Please try again.")
		return
	}

	roomNo = getRoomNo(roomType)

	fmt.Printf("Booking a %s room for %s (%s) for %d nights.\n", roomType, name, email, nights)
	fmt.Print("CheckIn date: ", checkIn, "\nCheckOut date: ", checkOut, "\n")
	var book = booking.Booking{
		Username: username, // Placeholder for username, can be replaced with actual user data
		Name:     name,
		Email:    email,
		RoomType: roomType,
		RoomNo:   roomNo,
		CheckIn:  checkIn,
		CheckOut: checkOut,
		Nights:   nights,
	}
	mu.Lock() // Lock the mutex to ensure thread safety
	// bookings = append(bookings, book)
	bookings[roomNo] = book
	mu.Unlock()
	fmt.Printf("Room booked successfully! Room No: %d\n", roomNo)
	go sendConfirmationEmail(name, email, roomNo, roomType, nights)
}

func checkAvailableRoom(roomType string) bool {
	switch roomType {
	case "Single":
		if singleRoomNoIdx < 5 {
			return true
		}
	case "Double":
		if doubleRoomNoIdx < 5 {
			return true
		}
	case "Suite":
		if suiteRoomNoIdx < 5 {
			return true
		}
	default:
		fmt.Println("Invalid room type. Please choose from Single, Double, or Suite.")
		return false
	}
	fmt.Println("No available rooms of this type.")
	return false
}

func getRoomNo(roomType string) int {
	// updateBookArrayfromJson()
	var roomNo int
	switch roomType {
	case "Single":
		singleRoomNoIdx++
		for i := 0; i < len(singleBookarray); i++ {
			if singleBookarray[i] == 0 {
				roomNo = singleRoomNo + i
				singleBookarray[i] = 1
				break
			}
		}
	case "Double":
		doubleRoomNoIdx++
		for i := 0; i < len(doubleBookarray); i++ {
			if doubleBookarray[i] == 0 {
				roomNo = doubleRoomNo + i
				doubleBookarray[i] = 1
				break
			}
		}
	case "Suite":
		suiteRoomNoIdx++
		for i := 0; i < len(suiteBookarray); i++ {
			if suiteBookarray[i] == 0 {
				roomNo = suiteRoomNo + i
				suiteBookarray[i] = 1
				break
			}
		}
	}
	return roomNo
}

func viewUserBookings(username string) {
	fmt.Println("Your Bookings:")
	var found bool = false
	for roomNo, booking := range bookings {
		if booking.Username == username {
			found = true
			fmt.Printf("Room No: %d, Name: %s, Email: %s, Room Type: %s, Nights: %d\n", roomNo, booking.Name, booking.Email, booking.RoomType, booking.Nights)
			fmt.Print("CheckIn date: ", booking.CheckIn, "\nCheckOut date: ", booking.CheckOut, "\n")
		}
	}
	if !found {
		fmt.Printf("You have no bookings.")
	}
}

func viewBookings() {
	fmt.Println("Current Bookings:")
	var totalRevenue int = 0
	for roomNo, booking := range bookings {
		fmt.Printf("Booked by Username %v.\nRoom No: %d, Name: %s, Email: %s, Room Type: %s, Nights: %d\n",booking.Username, roomNo, booking.Name, booking.Email, booking.RoomType, booking.Nights)
		fmt.Print("CheckIn date: ", booking.CheckIn, "\nCheckOut date: ", booking.CheckOut, "\n")
		if booking.RoomType == "Single" {
			totalRevenue = totalRevenue + (booking.Nights * 100)
		} else if booking.RoomType == "Double" {
			totalRevenue = totalRevenue + (booking.Nights * 200)
		} else if booking.RoomType == "Suite" {
			totalRevenue = totalRevenue + (booking.Nights * 300)
		} else {
			fmt.Println("Invalid room type in booking.")
		}
	}
	fmt.Printf("\nTotal Revenue from Bookings: $%d\n\n", totalRevenue)

	showAvailabilityCalendar()
}



func cancelBooking(username string) {
	var roomNo int
	fmt.Print("Enter the Room No to cancel booking: ")
	fmt.Scan(&roomNo)

	if bookings[roomNo].Username != username {
		fmt.Println("You can only cancel your own bookings.")
		return
	}
	mu.Lock() // Lock the mutex to ensure thread safety
	if bookings[roomNo].RoomNo == 0 {
		fmt.Println("No booking found for Room No:", roomNo)
		return
	}
	if bookings[roomNo].RoomType == "Single" {
		singleRoomNoIdx--
		singleBookarray[roomNo-singleRoomNo] = 0
	} else if bookings[roomNo].RoomType == "Double" {
		doubleRoomNoIdx--
		doubleBookarray[roomNo-doubleRoomNo] = 0
	} else if bookings[roomNo].RoomType == "Suite" {
		suiteRoomNoIdx--
		suiteBookarray[roomNo-suiteRoomNo] = 0
	} else {
		fmt.Println("Invalid room type for cancellation.")
		return
	}
	delete(bookings, roomNo)
	fmt.Println("Booking for Room No", roomNo, "has been cancelled successfully.")
	saveBookingsToFile("bookingData.json", bookings)
	mu.Unlock() // Unlock the mutex after modifying bookings
}

func sendConfirmationEmail(name, email string, roomNo int, roomType string, nights int) {
	time.Sleep(10 * time.Second) 
	fmt.Println("\n#########################################################################")
	fmt.Printf("Sending confirmation email to %s (%s) for booking Room No: %d (%s) for %d nights.\n", name, email, roomNo, roomType, nights)
	// Here you would implement the actual email sending logic
	// For now, we just print a message
	fmt.Println("Confirmation email sent successfully!")
	fmt.Println("#########################################################################")
	saveBookingsToFile("bookingData.json", bookings)
	
}

func saveBookingsToFile(filename string, bookings map[int]booking.Booking) error {
	data, err := json.MarshalIndent(bookings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func loadBookingsFromFile(filename string) (map[int]booking.Booking) {
	bookings := make(map[int]booking.Booking) // Initialize the map to hold bookings
	// bookings := make(map[int]booking.Booking)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return bookings // Return empty map if file doesn't exist
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &bookings)
	if err != nil {
		return nil
	}

	return bookings
}

func updateBookArrayfromJson() {
	for i := 0; i < len(singleBookarray); i++ {
		singleBookarray[i] = 0
	}
	for i := 0; i < len(doubleBookarray); i++ {
		doubleBookarray[i] = 0
	}
	for i := 0; i < len(suiteBookarray); i++ {
		suiteBookarray[i] = 0
	}

	for _, booking := range bookings {
		if booking.RoomType == "Single" {
			singleBookarray[booking.RoomNo-singleRoomNo] = 1
		} else if booking.RoomType == "Double" {
			doubleBookarray[booking.RoomNo-doubleRoomNo] = 1
		} else if booking.RoomType == "Suite" {
			suiteBookarray[booking.RoomNo-suiteRoomNo] = 1
		}
	}
}

func userLogin() string{
	
	var logindata = loadBookingsFromlogin("login.json")
	fmt.Println("Welcome to the Hotel Booking System")
	fmt.Println("1. Login")
	fmt.Println("2. Sign Up")
	fmt.Println("3. Exit")
	var choice int
	fmt.Print("Please enter your choice: ")
	fmt.Scan(&choice)
	
	if choice == 2 {
		var username, password string
		fmt.Println("Enter your details to sign up:")
		fmt.Print("Username: ")
		fmt.Scan(&username)
		fmt.Print("Password: ")
		fmt.Scan(&password)
		if _, exists := logindata[username]; exists {
			fmt.Println("Username already exists. Please choose a different username.")
			return "f"
		}

		// Hash the password before storing it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Error hashing password:", err)
			return "f"
		}
		// Store the hashed password in the logindata map
		password = string(hashedPassword) // Convert hashed password back to string for storage
		logindata[username] = password
		saveBookingsTologin("login.json", logindata)
		return username
	} else if choice == 3 {
		fmt.Println("Exiting the system. Thank you!")
		os.Exit(0) // Exit the program if user chooses to exit
	} else if choice == 1 {
		
		var username, password string
		fmt.Println("Please login to the Hotel Booking System")
		fmt.Print("Username: ")
		fmt.Scan(&username)
		fmt.Print("Password: ")
		fmt.Scan(&password)

		if username == "admin" && password == logindata["admin"]{
			fmt.Println("Admin login successful!")
			viewBookings() // Admin can view all bookings
			os.Exit(0) // Exit after viewing bookings
		}

		storedHash := logindata[username]
		err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
		if err == nil {
    		fmt.Println("Login successful!")
    	return username
		} else {
    		fmt.Println("Invalid username or password. Please try again.")
    		return "f"
		}		


	}
	return "f"
}
func saveBookingsTologin(filename string, logindata map[string]string) error {
	data, err := json.MarshalIndent(logindata, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
func loadBookingsFromlogin(filename string) (map[string]string) {
	logindata := make(map[string]string) // Initialize the map to hold login data
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return logindata // Return empty map if file doesn't exist
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &logindata)
	if err != nil {
		return nil
	}

	return logindata
}




func showAvailabilityCalendar() {
	// bookings := loadBookingsFromFile("bookingData.json")
	if len(bookings) == 0 {
		fmt.Println("No bookings available to generate calendar.")
		return
	}

	var earliest, latest time.Time

	// Step 1: Find earliest check-in and latest check-out
	for _, b := range bookings {
		if earliest.IsZero() || b.CheckIn.Before(earliest) {
			earliest = b.CheckIn
		}
		if latest.IsZero() || b.CheckOut.After(latest) {
			latest = b.CheckOut
		}
	}

	// Step 2: Generate date range
	dateRange := getDateRange(earliest, latest)

	// Step 3: Print calendar header
	fmt.Println("\nAvailability Calendar for Single Rooms:")
	fmt.Printf("Room No |")
	for _, d := range dateRange {
		fmt.Printf(" %s ", d.Format("02")) // Just print day-of-month
	}
	fmt.Println()

	// Step 4: Print row for each room
	for roomNo := 101; roomNo <= 105; roomNo++ {
		fmt.Printf("  %3d   |", roomNo)
		booking, exists := bookings[roomNo]
		for _, d := range dateRange {
			if exists && isBookedOnDate(booking, d) {
				fmt.Print("  B ")
			} else {
				fmt.Print("  F ")
			}
		}
		fmt.Println()
	}

	fmt.Println("\nAvailability Calendar for Double Rooms:")
	fmt.Printf("Room No |")
	for _, d := range dateRange {
		fmt.Printf(" %s ", d.Format("02")) // Just print day-of-month
	}
	fmt.Println()

	for roomNo := 201; roomNo <= 205; roomNo++ {
		fmt.Printf("  %3d   |", roomNo)
		booking, exists := bookings[roomNo]
		for _, d := range dateRange {
			if exists && isBookedOnDate(booking, d) {
				fmt.Print("  B ")
			} else {
				fmt.Print("  F ")
			}
		}
		fmt.Println()
	}

	fmt.Println("\nAvailability Calendar for Suite Rooms:")
	fmt.Printf("Room No |")
	for _, d := range dateRange {
		fmt.Printf(" %s ", d.Format("02")) // Just print day-of-month
	}
	fmt.Println()

	for roomNo := 301; roomNo <= 305; roomNo++ {
		fmt.Printf("  %3d   |", roomNo)
		booking, exists := bookings[roomNo]
		for _, d := range dateRange {
			if exists && isBookedOnDate(booking, d) {
				fmt.Print("  B ")
			} else {
				fmt.Print("  F ")
			}
		}
		fmt.Println()
	}
}

// Supporting functions below 👇

func getDateRange(start, end time.Time) []time.Time {
	dates := []time.Time{}
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}
	return dates
}

func isBookedOnDate(b booking.Booking, date time.Time) bool {
	date = date.Truncate(24 * time.Hour)
	start := b.CheckIn.Truncate(24 * time.Hour)
	end := b.CheckOut.Truncate(24 * time.Hour)
	return !date.Before(start) && date.Before(end)
}
