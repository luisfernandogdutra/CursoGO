package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	initTracing()

	http.HandleFunc("/cep", handleCEP)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Serviço A rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("service-a")
	_, span := tracer.Start(ctx, "handleCEP")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CEP string `json:"cep"`
	}
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &req); err != nil || len(req.CEP) != 8 {
		http.Error(w, `{"message":"invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	resp, err := http.Post("http://service-b:8081/cep", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Service B unavailable", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	//w.Write([]byte(resp.Body))
	json.NewEncoder(w).Encode(resp.Body)
}

func initTracing() {
	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("Erro ao inicializar Zipkin: %v", err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes()))
	otel.SetTracerProvider(tp)
}
