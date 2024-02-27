package db

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: os.Getenv("REDIS_PW"),
	DB:       0, // use default DB
})
