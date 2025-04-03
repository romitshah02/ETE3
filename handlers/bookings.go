package handlers

import (
	"fmt"
	"net/http"

	"ETE3/db"
	"ETE3/models"

	"github.com/gin-gonic/gin"
)

// BookSeats handles the booking of multiple seats for a show
func BookSeats(c *gin.Context) {
	// Use a fixed userID for testing
	userID := uint(1)

	var bookingRequest struct {
		ShowID uint     `json:"show_id" binding:"required"`
		Seats  []string `json:"seats" binding:"required,min=1"`
	}

	// Bind the JSON request data to the struct
	if err := c.ShouldBindJSON(&bookingRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Begin a transaction
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Fetch the show details using ShowID
	var show models.Show
	if err := tx.First(&show, bookingRequest.ShowID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
		return
	}

	// Get all seats for this show first, then filter in memory
	var allShowSeats []models.Seat
	if err := tx.Raw("SELECT * FROM seats WHERE show_id = ? AND deleted_at IS NULL",
		bookingRequest.ShowID).Scan(&allShowSeats).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seats"})
		return
	}

	// Create a map for easy lookup
	seatMap := make(map[string]models.Seat)
	for _, seat := range allShowSeats {
		key := fmt.Sprintf("%s%d", seat.Row, seat.Number)
		seatMap[key] = seat
	}

	// Validate each requested seat
	var seatsToBook []models.Seat
	for _, label := range bookingRequest.Seats {
		seat, exists := seatMap[label]
		if !exists {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Seat %s not found", label),
			})
			return
		}

		if seat.Status != models.Available {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("Seat %s is already booked", label),
			})
			return
		}

		seatsToBook = append(seatsToBook, seat)
	}

	// Update each seat status directly in the database
	for _, seat := range seatsToBook {
		// Use raw SQL to update status
		if err := tx.Exec("UPDATE seats SET status = ? WHERE id = ?",
			models.Booked, seat.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update seat status",
			})
			return
		}
	}

	// Create a new booking
	booking := models.Booking{
		UserID: userID,
		ShowID: bookingRequest.ShowID,
		Status: "confirmed",
	}

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create booking",
		})
		return
	}

	// Add seat associations using raw SQL to avoid GORM issues
	for _, seat := range seatsToBook {
		if err := tx.Exec("INSERT INTO booking_seats (booking_id, seat_id) VALUES (?, ?)",
			booking.ID, seat.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to associate seats with booking",
			})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to complete booking transaction",
		})
		return
	}

	// Calculate total price
	totalPrice := show.Price * float64(len(bookingRequest.Seats))

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":     "Booking confirmed",
		"booking_id":  booking.ID,
		"show_id":     show.ID,
		"seats":       bookingRequest.Seats,
		"total_price": totalPrice,
		"status":      booking.Status,
	})
}
