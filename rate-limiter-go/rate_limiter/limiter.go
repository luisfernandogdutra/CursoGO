package rate_limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Limiter struct {
	client  *redis.Client
	timeout time.Duration
}

func NewLimiter(redisClient *redis.Client, timeout time.Duration) *Limiter {
	return &Limiter{
		client:  redisClient,
		timeout: timeout,
	}
}

func (l *Limiter) IsAllowed(key string, limit int) bool {
	ctx := context.Background()
	pipe := l.client.TxPipeline()

	// Increment the count of requests
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, l.timeout)

	_, err := pipe.Exec(ctx)
	if err != nil {
		fmt.Println("Error executing Redis pipeline:", err)
		return false
	}

	// If count exceeds limit, reject the request
	return incr.Val() <= int64(limit)
}

func (l *Limiter) KeyForIP(ip string) string {
	return fmt.Sprintf("ip:%s", ip)
}

func (l *Limiter) KeyForToken(token string) string {
	return fmt.Sprintf("token:%s", token)
}
