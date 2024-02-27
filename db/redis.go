package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func getPw() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pw := os.Getenv("REDIS_PW")
	if pw == "" {
		log.Fatal("No Redis PW set")
	}
	return pw
}

var Redis *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: getPw(),
	DB:       0, // use default DB
})
