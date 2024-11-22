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
		log.Fatal("Error ao carregar .env")
	}

	limit := 10
	if envLimit := os.Getenv("RATE_LIMIT"); envLimit != "" {
		limit, err = strconv.Atoi(envLimit)
		if err != nil {
			fmt.Println("Erro duranter a conversao")
			return 0, nil
		}
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	client := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	return limit, client
}
