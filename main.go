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
		{Title: "Inception", Duration: 148, Photo: "https://imgs.search.brave.com/ewibJwgPR-UwtznSPey5xuPBdlBdQxnzqS2L8aOMFbc/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly9pcnMu/d3d3Lndhcm5lcmJy/b3MuY29tL2tleWFy/dC1qcGVnL2luY2Vw/dGlvbl9rZXlhcnQu/anBn"},
		{Title: "The Dark Knight", Duration: 152, Photo: "https://imgs.search.brave.com/8JdUxOLqLfuVw3GmzUtMWijBb-W7BgK0j8VgWNtgMBQ/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly9pbWFn/ZXMtbmEuc3NsLWlt/YWdlcy1hbWF6b24u/Y29tL2ltYWdlcy9J/LzgxK1lud1J1Vy1M/LmpwZw"},
		{Title: "Interstellar", Duration: 169, Photo: "https://imgs.search.brave.com/oFOlPK2YLX9auJAjul2GCtXsyiZA7_AqLl-MNgZPwLs/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly9pbWFn/ZXMtbmEuc3NsLWlt/YWdlcy1hbWF6b24u/Y29tL2ltYWdlcy9J/LzgxWWxyUVo2WjRT/LmpwZw"},
		{Title: "The Matrix", Duration: 136, Photo: "https://imgs.search.brave.com/B_dNIhoDAHtAIIPC-Drqkyruy42iTD1SzYvF3oA_QH0/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly9tLm1l/ZGlhLWFtYXpvbi5j/b20vaW1hZ2VzL00v/TVY1Qk4yTm1OMlZo/TVRRdE1ETmlPUzAw/TkRsaExUbGlNamd0/T0RFMlpUWTBPRFF5/TkRSaFhrRXlYa0Zx/Y0djQC5qcGc"},
		{Title: "Avatar", Duration: 162, Photo: "https://imgs.search.brave.com/o65ljDRFOI0NqiWW59I1Kbyt0nxR6O-YjiGaMqgsNrs/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly9tLm1l/ZGlhLWFtYXpvbi5j/b20vaW1hZ2VzL00v/TVY1Qk1ERXpNbVF3/WmpjdFpXVTJNeTAw/TVdObExXRTBOakl0/TURKbFlUUmxOR0pp/WmpjeVhrRXlYa0Zx/Y0djQC5qcGc"},
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
