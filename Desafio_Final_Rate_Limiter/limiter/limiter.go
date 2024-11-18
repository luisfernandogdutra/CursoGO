package limiter

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type RateLimiter struct {
	redisClient *redis.Client
	ipLimit     int
	tokenLimit  int
	blockTime   int
}

func NewRateLimiter() *RateLimiter {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	redisURL := os.Getenv("URL")
	ipLimit, _ := strconv.Atoi(os.Getenv("LIMIT_IP"))
	tokenLimit, _ := strconv.Atoi(os.Getenv("LIMIT_TOKEN"))
	blockTime, _ := strconv.Atoi(os.Getenv("BLOCK_TIME"))

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &RateLimiter{
		redisClient: rdb,
		ipLimit:     ipLimit,
		tokenLimit:  tokenLimit,
		blockTime:   blockTime,
	}
}

func (rl *RateLimiter) IsRateLimited(ip, token string) (bool, error) {
	var limit int
	var key string

	if token != "" {
		key = "token:" + token
		limit = rl.tokenLimit
	} else {
		key = "ip:" + ip
		limit = rl.ipLimit
	}

	ctx := context.Background()
	currentCount, err := rl.redisClient.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("ERRO AO BUSCAR DO REDIS: %w", err)
	}

	if currentCount >= limit {
		return true, nil
	}

	ttl := time.Duration(rl.blockTime) * time.Second
	err = rl.redisClient.Incr(ctx, key).Err()
	if err != nil {
		return false, fmt.Errorf("ERRO AO INCREMENTAR O CONTADOR: %w", err)
	}

	rl.redisClient.Expire(ctx, key, ttl)
	return false, nil
}
