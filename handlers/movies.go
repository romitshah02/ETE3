package handlers

import (
	"ETE3/db"
	"ETE3/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.DB.Create(&movie)
	c.JSON(200, gin.H{"message": "Movie added successfully", "movie_id": movie.ID})
}

func AddShowHandler(c *gin.Context) {
	var show models.Show

	// Validate request body
	if err := c.ShouldBindJSON(&show); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the show to the database
	if err := db.DB.Create(&show).Error; err != nil {
		log.Printf("Error adding show: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create show"})
		return
	}

	// Generate seats for this show (A1 to J15)
	rows := "ABCDEFGHIJ" // 10 rows (A to J)
	seats := []models.Seat{}

	for _, row := range rows {
		for num := 1; num <= 15; num++ {
			seats = append(seats, models.Seat{
				ShowID: show.ID,
				Row:    string(row),
				Number: num,
				Status: models.Available,
			})
		}
	}

	// Insert seats into the DB
	if err := db.DB.Create(&seats).Error; err != nil {
		log.Printf("Error adding seats for show %d: %v", show.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seats"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Show created successfully",
		"show":    show,
	})
}

func GetAllMovies(c *gin.Context) {
	var movies []models.Movie

	// Fetch all movies from the database
	if err := db.DB.Find(&movies).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to fetch movies"})
		return
	}

	// Return the list of movies as JSON
	c.JSON(200, movies)
}

func GetShowsByMovie(c *gin.Context) {
	movieID := c.Param("movie_id") // Retrieve the movie ID from the URL path

	// Query for all shows related to the given movie ID
	var shows []models.Show
	if err := db.DB.Where("movie_id = ?", movieID).Find(&shows).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to fetch shows"})
		return
	}

	// If no shows are found
	if len(shows) == 0 {
		c.JSON(404, gin.H{"message": "No shows found for this movie"})
		return
	}

	// Return the list of shows as JSON
	c.JSON(200, gin.H{"shows": shows})
}

func GetAvailableSeatsHandler(c *gin.Context) {
	var seats []models.Seat
	showID := c.Param("show_id") // Extract show_id from URL

	// Fetch only available seats for the given show
	if err := db.DB.Where("show_id = ? AND status = ?", showID, models.Available).Find(&seats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available seats"})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"show_id":         showID,
		"seats":           seats,
		"total_available": len(seats),
	})
}
