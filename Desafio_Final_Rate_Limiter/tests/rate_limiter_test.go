package tests

import (
	"testing"
	"time"

	"github.com/luisfernandogdutra/CursoGO/Desafio_Final_Rate_Limiter/limiter"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	rateLimiter := limiter.NewRateLimiter()

	t.Run("Should not rate limit on first request by IP", func(t *testing.T) {
		ip := "192.168.0.1"
		limited, err := rateLimiter.IsRateLimited(ip, "")
		assert.NoError(t, err, "Error should be nil on first request by IP")
		assert.False(t, limited, "First request should not be rate limited")
	})

	t.Run("Should not rate limit on first request by Token", func(t *testing.T) {
		token := "tokenteste123"
		limited, err := rateLimiter.IsRateLimited("", token)
		assert.NoError(t, err, "Error should be nil on first request by Token")
		assert.False(t, limited, "First request should not be rate limited")
	})

	t.Run("Should rate limit after exceeding limit by IP", func(t *testing.T) {
		ip := "192.168.0.2"
		for i := 0; i < 5; i++ {
			limited, err := rateLimiter.IsRateLimited(ip, "")
			assert.NoError(t, err, "Error should be nil on requests by IP")
			assert.False(t, limited, "Requests within limit should not be rate limited")
		}
		// Exceed the limit
		limited, err := rateLimiter.IsRateLimited(ip, "")
		assert.NoError(t, err, "Error should be nil on rate-limited request by IP")
		assert.True(t, limited, "Should rate limit after exceeding request limit by IP")
	})

	t.Run("Should rate limit after exceeding limit by Token", func(t *testing.T) {
		token := "tokenteste456"
		for i := 0; i < 5; i++ {
			limited, err := rateLimiter.IsRateLimited("", token)
			assert.NoError(t, err, "Error should be nil on requests by Token")
			assert.False(t, limited, "Requests within limit should not be rate limited")
		}
		// Exceed the limit
		limited, err := rateLimiter.IsRateLimited("", token)
		assert.NoError(t, err, "Error should be nil on rate-limited request by Token")
		assert.True(t, limited, "Should rate limit after exceeding request limit by Token")
	})

	t.Run("Should reset rate limiting after time passes", func(t *testing.T) {
		ip := "192.168.0.3"
		token := "tokenteste789"

		// Exceed the limit
		for i := 0; i < 5; i++ {
			rateLimiter.IsRateLimited(ip, token)
		}

		limited, _ := rateLimiter.IsRateLimited(ip, token)
		assert.True(t, limited, "Should rate limit after exceeding the limit")

		// Wait for the rate limit to reset (adjust duration as per implementation)
		time.Sleep(2 * time.Second) // Ajuste o tempo conforme o reset configurado no RateLimiter

		limited, err := rateLimiter.IsRateLimited(ip, token)
		assert.NoError(t, err, "Error should be nil after rate limit reset")
		assert.False(t, limited, "Rate limit should reset after the specified time")
	})
}
