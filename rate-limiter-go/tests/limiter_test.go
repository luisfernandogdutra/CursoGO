package tests

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/luisfernandogdutra/rate-limiter-go/rate_limiter"
	"github.com/stretchr/testify/assert"
)

func TestLimiter_IsAllowed(t *testing.T) {
	// Criar uma instância de cliente Redis em memória (utilizando o Redis default)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Pode ser configurado para usar Redis em memória no CI
	})
	limiter := rate_limiter.NewLimiter(client, 2*time.Second)

	t.Run("Test Allow Request Below Limit", func(t *testing.T) {
		// Limite de 2 requisições por chave
		key := "test_key"
		allowed := limiter.IsAllowed(key, 2)
		assert.True(t, allowed, "Request should be allowed")
	})

	t.Run("Test Block Request Above Limit", func(t *testing.T) {
		key := "test_key"
		limiter.IsAllowed(key, 2)            // Primeira requisição (permitida)
		limiter.IsAllowed(key, 2)            // Segunda requisição (permitida)
		blocked := limiter.IsAllowed(key, 2) // Terceira requisição (deve ser bloqueada)
		assert.False(t, blocked, "Request should be blocked")
	})
}

func TestLimiter_Expiration(t *testing.T) {
	// Criar uma instância de cliente Redis em memória
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Pode ser configurado para usar Redis em memória no CI
	})
	limiter := rate_limiter.NewLimiter(client, 2*time.Second)

	t.Run("Test Expiration of Keys", func(t *testing.T) {
		key := "test_key"
		// Primeiro, vamos permitir 2 requisições
		limiter.IsAllowed(key, 2)
		limiter.IsAllowed(key, 2)

		// A terceira requisição será bloqueada
		blocked := limiter.IsAllowed(key, 2)
		assert.False(t, blocked, "Request should be blocked")

		// Espera o tempo de expiração (2 segundos) e tenta novamente
		time.Sleep(2 * time.Second)

		// Agora, a chave deve ser desbloqueada e a requisição permitida
		allowed := limiter.IsAllowed(key, 2)
		assert.True(t, allowed, "Request should be allowed after expiration")
	})
}
