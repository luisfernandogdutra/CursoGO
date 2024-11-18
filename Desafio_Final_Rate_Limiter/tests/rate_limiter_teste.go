package limiter

import (
	"testing"

	"github.com/luisfernandogdutra/CursoGO/Desafio_Final_Rate_Limiter/limiter"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	rateLimiter := limiter.NewRateLimiter()

	ip := "192.168.0.1"
	token := "tokenteste123"

	limited, err := rateLimiter.IsRateLimited(ip, "")
	assert.NoError(t, err)
	assert.False(t, limited)

	limited, err = rateLimiter.IsRateLimited("", token)
	assert.NoError(t, err)
	assert.False(t, limited)
}
