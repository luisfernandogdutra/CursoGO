package main

import (
	"log"
	"net/http"

	"github.com/luisfernandogdutra/desafio-rate-limiter/rate_limiter"
)

func main() {
	limit, redisClient := rate_limiter.LoadConfig()
	http.HandleFunc("/", rate_limiter.RateLimitMiddleware(redisClient, limit))

	log.Println("Iniciando server: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error ao iniciar server: ", err)
	}
}
