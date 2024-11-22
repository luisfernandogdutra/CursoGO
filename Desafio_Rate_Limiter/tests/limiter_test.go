package tests

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/luisfernandogdutra/desafio-rate-limiter/rate_limiter"
	"github.com/stretchr/testify/assert"
)

func TestLimiter_IsAllowed(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	limiter := rate_limiter.NewLimiter(client, 2*time.Second)

	t.Run("Teste Requisicao Abaixo do Limite", func(t *testing.T) {
		key := "test_key"
		allowed := limiter.IsAllowed(key, 2)
		assert.True(t, allowed, "Requisicao dever ser permitida")
	})

	t.Run("Teste Bloqueio de Requisicao Acima do Limite", func(t *testing.T) {
		key := "test_key"
		limiter.IsAllowed(key, 2)
		limiter.IsAllowed(key, 2)
		blocked := limiter.IsAllowed(key, 2)
		assert.False(t, blocked, "Requisicao deve ser bloqueada")
	})
}

func TestLimiter_Expiration(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	limiter := rate_limiter.NewLimiter(client, 2*time.Second)

	t.Run("Teste Expiracao da Chave", func(t *testing.T) {
		key := "test_key"
		limiter.IsAllowed(key, 2)
		limiter.IsAllowed(key, 2)

		blocked := limiter.IsAllowed(key, 2)
		assert.False(t, blocked, "Requisicao deve ser bloqueada")

		time.Sleep(2 * time.Second)
		allowed := limiter.IsAllowed(key, 2)
		assert.True(t, allowed, "Requisicao deve ser permitida depois do bloqueio")
	})
}
