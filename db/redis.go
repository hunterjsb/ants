package db

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func getPw() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pw := os.Getenv("REDIS_PW")
	return pw
}

func getDb() int {
	rdb := os.Getenv("REDIS_DB")
	rdbNum, err := strconv.Atoi(rdb)
	if err != nil {
		log.Println("Invalid Redis DB or DB not supplied - using 0")
		return 0
	}
	return rdbNum
}

var Redis *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: getPw(),
	DB:       getDb(),
})
