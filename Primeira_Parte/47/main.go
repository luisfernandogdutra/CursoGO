package main

import (
	"log"
	"net/http"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered panic: %v\n", r)
				//debug.PrintStack()
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello World\n"))
	})

	mux.HandleFunc("/panic", func(writer http.ResponseWriter, request *http.Request) {
		panic("panic")
	})

	log.Println("Listening on", ":3000")
	if err := http.ListenAndServe(":3000", recoverMiddleware(mux)); err != nil {
		log.Fatalf("Could not listen on %s: %v\n", ":3000", err)
	}
}
