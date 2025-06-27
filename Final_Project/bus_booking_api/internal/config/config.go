package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type Config struct {
	Port                string
	MongoURI            string
	MongoDBName         string
	JWTSecretKey        string
	JWTExpiration       time.Duration
}

var AppConfig Config
var DB *mongo.Database

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig.Port = getEnv("PORT", "8080")
	AppConfig.MongoURI = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	AppConfig.MongoDBName = getEnv("MONGODB_DATABASE_NAME", "bus_booking_db")
	AppConfig.JWTSecretKey = getEnv("JWT_SECRET_KEY", "default_secret_key_please_change")

	jwtExpHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "72"))
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRATION_HOURS: %v", err)
	}
	AppConfig.JWTExpiration = time.Duration(jwtExpHours) * time.Hour

	if AppConfig.JWTSecretKey == "default_secret_key_please_change" {
		log.Println("WARNING: JWT_SECRET_KEY is set to default. Please change it for production.")
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func ConnectDB() {
	clientOptions := options.Client().ApplyURI(AppConfig.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	DB = client.Database(AppConfig.MongoDBName)
}

func CloseDB(client *mongo.Client) {
	if client == nil {
		return
	}
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
	log.Println("Connection to MongoDB closed.")
}

func GetCollection(collectionName string) *mongo.Collection {
    if DB == nil {
        log.Fatal("Database not initialized. Call ConnectDB first.")
    }
	return DB.Collection(collectionName)
}
