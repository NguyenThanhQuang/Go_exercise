package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func ConnectRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
		log.Println("REDIS_ADDR not set in .env, using default:", redisAddr)
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil || redisDBStr == "" {
		redisDB = 0
		log.Println("REDIS_DB not set or invalid in .env, using default:", redisDB)
	}

	redisUseTLSStr := os.Getenv("REDIS_USE_TLS")
	redisUseTLS, _ := strconv.ParseBool(redisUseTLSStr) 

	redisOptions := &redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	}

	if redisUseTLS {
		log.Println("Attempting to connect to Redis with TLS.")
		redisOptions.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	} else {
		log.Println("Attempting to connect to Redis without TLS.")
	}


	RDB = redis.NewClient(redisOptions)

	pong, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis successfully! Ping response: %s\n", pong)
}