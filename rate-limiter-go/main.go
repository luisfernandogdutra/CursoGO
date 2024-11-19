package main

import (
	"log"
	"net/http"

	"github.com/luisfernandogdutra/rate-limiter-go/rate_limiter"
)

func main() {
	limit, redisClient := rate_limiter.LoadConfig()
	http.HandleFunc("/", rate_limiter.RateLimitMiddleware(redisClient, limit))

	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
