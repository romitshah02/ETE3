package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	log.Default().Println("Initializing the connection to database.")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",
		"Romit@744812",
		"127.0.0.1",
		"go_ete")

	db, err := gorm.Open(mysql.Open(dataSourceName))
	if err != nil {
		log.Fatalln("Database Initialization failed. ", err)
	}
	log.Println("Connected to database. Host: ", os.Getenv("DB_HOST"), ". Database:", os.Getenv("DB_NAME"))
	DB = db
}

func CloseDB() {
	sqlDB, _ := DB.DB()
	log.Default().Println("Closing database connection")
	sqlDB.Close()
	log.Default().Println("Connection closed")
}
