package models

import (
	"time"

	"gorm.io/gorm"
)

type bookingStatus int

const (
	Available bookingStatus = iota
	Booked
)

func (s bookingStatus) String() string {
	return [...]string{"available", "booked"}[s]
}

type Seat struct {
	gorm.Model
	ShowID uint          `json:"show_id"`
	Row    string        `json:"row"`    // e.g., "A"
	Number int           `json:"number"` // e.g., 1-15
	Status bookingStatus `json:"status" gorm:"default:0"`
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}

type Movie struct {
	gorm.Model
	Title    string `json:"title"`
	Duration int    `json:"duration"` // in minutes
	Photo    string `json:"photo"`    // Store the photo URL or file path
}

type Show struct {
	gorm.Model
	MovieID uint      `json:"movie_id"`
	Time    time.Time `json:"time"` // Use time.Time for handling date and time
	Price   float64   `json:"price"`
}

type Booking struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	ShowID uint   `json:"show_id"`
	Seats  []Seat `json:"seats" gorm:"many2many:booking_seats;"`
	Status string `json:"status"` // confirmed, cancelled
}
