package rate_limiter

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func RateLimitMiddleware(redisClient *redis.Client, limit int) http.HandlerFunc {
	limiter := NewLimiter(redisClient, time.Second)

	return func(w http.ResponseWriter, r *http.Request) {
		var key string
		token := r.Header.Get("API_KEY")
		if token != "" {
			key = limiter.KeyForToken(token)
		} else {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			key = limiter.KeyForIP(ip)
		}

		if !limiter.IsAllowed(key, limit) {
			http.Error(w, "Server alcancou o numero maximo de requisicoes", http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Requisicao permitida"))
	}
}
