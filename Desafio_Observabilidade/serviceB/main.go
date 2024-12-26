// main.go - Serviço B
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	initTracing()

	http.HandleFunc("/cep", handleCEP)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Serviço B rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "handleCEP")
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

	city, err := getCityFromCEP(ctx, req.CEP)
	if err != nil {
		http.Error(w, `{"message":"can not find zipcode"}`, http.StatusNotFound)
		return
	}
	tempC, err := getTemperature(ctx, city)
	if err != nil {
		http.Error(w, "Failed to fetch temperature", http.StatusInternalServerError)
		return
	}
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	resp := struct {
		City  string  `json:"city"`
		TempC float64 `json:"temp_C"`
		TempF float64 `json:"temp_F"`
		TempK float64 `json:"temp_K"`
	}{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func getCityFromCEP(ctx context.Context, cep string) (string, error) {
	tracer := otel.Tracer("service-b")
	_, span := tracer.Start(ctx, "getCityFromCEP")
	defer span.End()

	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CEP not found")
	}

	var result struct {
		Localidade string `json:"localidade"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Localidade, nil
}

func getTemperature(ctx context.Context, city string) (float64, error) {
	tracer := otel.Tracer("service-b")
	_, span := tracer.Start(ctx, "getTemperature")
	defer span.End()

	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=bad4fe6d4148402daa712525242412&q=%s", encodedCity)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to fetch weather data")
	}

	var result struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Current.TempC, nil
}

func initTracing() {
	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("Erro ao inicializar Zipkin: %v", err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes()),
	)
	otel.SetTracerProvider(tp)
}
