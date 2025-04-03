package handlers

import (
	"ETE3/db"
	"ETE3/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BookSeats(c *gin.Context) {
	var bookingRequest struct {
		ShowID uint   `json:"show_id"`
		Seats  []uint `json:"seats"` // List of seat IDs the user wants to book
	}

	// Bind JSON request data
	if err := c.ShouldBindJSON(&bookingRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get user ID from token stored in context
	// userID, exists := c.Get("id")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
	// 	return
	// }

	userID, _ := strconv.ParseUint(c.Param("id"), 10, 0)

	// Fetch the show details
	var show models.Show
	if err := db.DB.First(&show, bookingRequest.ShowID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
		return
	}

	// Fetch the requested seats
	var seats []models.Seat
	if err := db.DB.Where("id IN ?", bookingRequest.Seats).Find(&seats).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Some seats not found"})
		return
	}

	// Check if all selected seats are available
	for _, seat := range seats {
		if seat.Status != models.Available {
			c.JSON(http.StatusConflict, gin.H{"error": "One or more seats are already booked"})
			return
		}
	}

	// Mark seats as booked
	for i := range seats {
		seats[i].Status = models.Booked
		if err := db.DB.Save(&seats[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seat status"})
			return
		}
	}

	// Create a booking record
	booking := models.Booking{
		UserID: uint(userID),
		ShowID: bookingRequest.ShowID,
		Seats:  seats,
		Status: "confirmed",
	}
	if err := db.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":    "Booking confirmed",
		"booking_id": booking.ID,
		"seats":      seats,
	})
}
