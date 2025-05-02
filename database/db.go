package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	dbhost := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	dbport := os.Getenv("DB_PORT")
	dbpass := os.Getenv("DB_PASSWORD")
	dbuser := os.Getenv("DB_USER")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbhost, dbuser, dbpass, dbname, dbport)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	DB = db
}
