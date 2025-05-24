package main

import (
	"log"
	"os"
	"social_media_server/config"

	"social_media_server/routes"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "social_media_server/docs" // This is for Swagger documentation
)

// @title Social Media API
// @version 1.0
// @description This is a sample server for a social media application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading, relying on environment variables")
	}

	config.ConnectDB()
	// config.ConnectRedis()

	router := routes.SetupRouter() 

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", port) 
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}