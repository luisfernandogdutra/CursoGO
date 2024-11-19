package rate_limiter

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func LoadConfig() (int, *redis.Client) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load rate limit from .env file
	limit := 10 // Default limit
	if envLimit := os.Getenv("RATE_LIMIT"); envLimit != "" {
		limit, err = strconv.Atoi(envLimit)

		if err != nil {
			fmt.Println("Error during conversion")
			return 0, nil
		}
	}

	// Setup Redis client
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	client := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	return limit, client
}
