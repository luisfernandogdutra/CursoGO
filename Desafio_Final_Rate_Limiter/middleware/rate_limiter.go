package middleware

import (
	"log"
	"net/http"

	"github.com/luisfernandogdutra/CursoGO/Desafio_Final_Rate_Limiter/limiter"
)

func RateLimiterMiddleware(rateLimiter *limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			token := r.Header.Get("API_KEY")

			isLimited, err := rateLimiter.IsRateLimited(ip, token)
			if err != nil {
				log.Println("ERRO AO CHECKAR LIMITE:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if isLimited {
				http.Error(w, "ALCANÇOU O LIMITE MÁXIMO DE REQUISIÇÕES", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
