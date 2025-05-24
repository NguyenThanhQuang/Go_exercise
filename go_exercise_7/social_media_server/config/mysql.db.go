package config

import (
	"fmt"
	"log"
	"os"
	"social_media_server/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB 

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("MYSQL_URL")
	if dsn == "" {
		log.Fatal("MYSQL_URL environment variable not set")
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connection successful!")

	err = database.AutoMigrate(&models.Post{}, &models.Comment{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
	fmt.Println("Database migration successful!")

	DB = database
}