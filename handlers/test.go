package handlers

import (
	"ETE3/db"
	"ETE3/models"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup an in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Movie{}, &models.Show{}, &models.Seat{})
	return db
}

// Test AddMovie Handler
func TestAddMovie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB()
	db.DB = testDB // Use test database

	router := gin.Default()
	router.POST("/movie", AddMovie)

	// Create a JSON request
	movieData := `{"title":"Test Movie"}`
	req, _ := http.NewRequest(http.MethodPost, "/movie", bytes.NewBuffer([]byte(movieData)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Movie added successfully")
}

// Test AddShowHandler
func TestAddShowHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB()
	db.DB = testDB // Use test database

	router := gin.Default()
	router.POST("/show", AddShowHandler)

	showData := `{"movie_id":1}`
	req, _ := http.NewRequest(http.MethodPost, "/show", bytes.NewBuffer([]byte(showData)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Show created successfully")
}

// Test GetAllMovies
func TestGetAllMovies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB()
	db.DB = testDB // Use test database

	// Insert test data
	testDB.Create(&models.Movie{Title: "Test Movie 1"})
	testDB.Create(&models.Movie{Title: "Test Movie 2"})

	router := gin.Default()
	router.GET("/movies", GetAllMovies)

	req, _ := http.NewRequest(http.MethodGet, "/movies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Movie 1")
	assert.Contains(t, w.Body.String(), "Test Movie 2")
}

// Test GetShowsByMovie
func TestGetShowsByMovie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB()
	db.DB = testDB // Use test database

	// Insert test data
	testDB.Create(&models.Movie{Title: "Test Movie"})
	testDB.Create(&models.Show{MovieID: 1})

	router := gin.Default()
	router.GET("/movie/:movie_id/shows", GetShowsByMovie)

	req, _ := http.NewRequest(http.MethodGet, "/movie/1/shows", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "shows")
}

// Test GetAvailableSeatsHandler
func TestGetAvailableSeatsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB()
	db.DB = testDB // Use test database

	// Insert test data
	testDB.Create(&models.Show{MovieID: 1})
	testDB.Create(&models.Seat{ShowID: 1, Row: "A", Number: 1, Status: models.Available})
	testDB.Create(&models.Seat{ShowID: 1, Row: "A", Number: 2, Status: models.Available})

	router := gin.Default()
	router.GET("/show/:show_id/seats", GetAvailableSeatsHandler)

	req, _ := http.NewRequest(http.MethodGet, "/show/1/seats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "seats")
}
