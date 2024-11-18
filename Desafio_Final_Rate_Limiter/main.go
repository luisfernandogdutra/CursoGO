package main

import (
	"log"
	"net/http"

	"github.com/luisfernandogdutra/CursoGO/Desafio_Final_Rate_Limiter/limiter"
	"github.com/luisfernandogdutra/CursoGO/Desafio_Final_Rate_Limiter/middleware"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	rateLimiter := limiter.NewRateLimiter()

	r.Use(middleware.RateLimiterMiddleware(rateLimiter))

	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Requisição feita com sucesso"))
	})

	log.Println("Server iniciando em: 8080")
	http.ListenAndServe(":8080", r)
}
