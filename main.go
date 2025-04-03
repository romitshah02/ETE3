package main

import (
	"ETE3/db"
	"ETE3/handlers"
	"ETE3/models"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	log.Default().Println("Starting server...")

	// Load environment variables
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Default().Println("Env file not found. Please check if you are using localhost. Err: ", envErr)
	}

	log.Default().Println("Initializing database...")
	db.InitDB()
	defer db.CloseDB()

	// Drop existing tables before migrating (ensure the tables are fresh)
	db.DB.Migrator().DropTable(&models.Movie{})
	db.DB.Migrator().DropTable(&models.User{})
	db.DB.Migrator().DropTable(&models.Show{})
	db.DB.Migrator().DropTable(&models.Seat{})
	db.DB.Migrator().DropTable(&models.Booking{})

	// AutoMigrate ensures that the schema matches the models
	db.DB.AutoMigrate(&models.User{})
	db.DB.AutoMigrate(&models.Movie{})
	db.DB.AutoMigrate(&models.Seat{})
	db.DB.AutoMigrate(&models.Booking{})
	db.DB.AutoMigrate(&models.Show{})

	// Seed movies and shows
	SeedMoviesAndShows()

	// Setup router and run the server
	r := handlers.SetupRouter()
	r.Run(":5000")
}

func SeedMoviesAndShows() {
	// List of movies to add
	movies := []models.Movie{
		{Title: "Inception", Duration: 148, Photo: "https://example.com/photos/inception.jpg"},
		{Title: "The Dark Knight", Duration: 152, Photo: "https://example.com/photos/dark_knight.jpg"},
		{Title: "Interstellar", Duration: 169, Photo: "https://example.com/photos/interstellar.jpg"},
		{Title: "The Matrix", Duration: 136, Photo: "https://example.com/photos/matrix.jpg"},
		{Title: "Avatar", Duration: 162, Photo: "https://example.com/photos/avatar.jpg"},
	}

	// Add movies to the database
	for _, movie := range movies {
		if err := db.DB.Create(&movie).Error; err != nil {
			log.Printf("❌ Error adding movie: %v", err)
			continue
		}
		log.Printf("✅ Movie added: %s", movie.Title)

		// Add shows for each movie
		addShowsForMovie(movie.ID)
	}
}

// addShowsForMovie adds shows and seats for a given movie
func addShowsForMovie(movieID uint) {
	shows := []models.Show{
		{MovieID: movieID, Time: time.Date(2025, time.April, 7, 10, 00, 0, 0, time.UTC), Price: 12.50},
		{MovieID: movieID, Time: time.Date(2025, time.April, 7, 14, 30, 0, 0, time.UTC), Price: 15.00},
		{MovieID: movieID, Time: time.Date(2025, time.April, 7, 18, 30, 0, 0, time.UTC), Price: 17.00},
		{MovieID: movieID, Time: time.Date(2025, time.April, 7, 22, 00, 0, 0, time.UTC), Price: 20.00},
	}

	// Add the shows for the movie
	for _, show := range shows {
		if err := db.DB.Create(&show).Error; err != nil {
			log.Printf("❌ Error adding show for movie %d: %v", movieID, err)
			continue
		}

		log.Printf("✅ Show added for movie %d at %s", movieID, show.Time.Format("15:04"))

		// Add seats for this show
		addSeatsForShow(show.ID)
	}
}

// addSeatsForShow adds 10 rows with 15 seats each for a show
func addSeatsForShow(showID uint) {
	var seats []models.Seat
	rows := "ABCDEFGHIJ" // 10 rows (A to J)

	// Generate seats for each row (A-J) and seat number (1-15)
	for _, row := range rows {
		for seatNum := 1; seatNum <= 15; seatNum++ {
			// Create a seat model with row and seat number
			seats = append(seats, models.Seat{
				ShowID: showID,
				Row:    string(row),      // Row is a string (A, B, C, ...)
				Number: seatNum,          // Seat number is an integer (1-15)
				Status: models.Available, // Initial status is Available
			})
		}
	}

	// Bulk insert seats into the database
	if err := db.DB.Create(&seats).Error; err != nil {
		log.Printf("❌ Error adding seats for show %d: %v", showID, err)
	} else {
		log.Printf("✅ Seats added for show %d", showID)
	}
}
