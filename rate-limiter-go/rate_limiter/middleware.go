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

		// Prefer Token over IP for rate limiting
		if token != "" {
			key = limiter.KeyForToken(token)
		} else {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			key = limiter.KeyForIP(ip)
		}

		// Check if request is allowed
		if !limiter.IsAllowed(key, limit) {
			http.Error(w, "You have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		// Proceed if allowed
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request allowed"))
	}
}
