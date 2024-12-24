package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestInvalidCEPFormat(t *testing.T) {
	invalidCEPs := []string{
		"1234-56",
		"123456789",
		"abcde-fgh",
	}

	for _, cep := range invalidCEPs {
		t.Run(fmt.Sprintf("CEP inválido: %s", cep), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/weather?cep=%s", cep))
			if err != nil {
				t.Fatalf("erro na requisição: %v", err)
			}
			if resp.StatusCode != http.StatusUnprocessableEntity {
				t.Errorf("esperado 422, mas obteve %v", resp.StatusCode)
			}
		})
	}
}

func TestCEPNotFound(t *testing.T) {
	cep := "00000-000"
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/weather?cep=%s", cep))
	if err != nil {
		t.Fatalf("erro na requisição: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("esperado 404, mas obteve %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("erro ao ler o corpo da resposta: %v", err)
	}

	expectedMessage := `{"message": "cidade não encontrada para o CEP"}`
	if strings.TrimSpace(string(body)) != expectedMessage {
		t.Errorf("esperado %s, mas obteve %s", expectedMessage, string(body))
	}
}

func TestSuccessfulWeatherQuery(t *testing.T) {
	cep := "01001-000"
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/weather?cep=%s", cep))
	if err != nil {
		t.Fatalf("erro na requisição: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado 200, mas obteve %v", resp.StatusCode)
	}

	var result map[string]float32
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("erro ao decodificar a resposta: %v", err)
	}

	if tempC, ok := result["temp_C"]; !ok || tempC == 0 {
		t.Errorf("esperado temp_C no corpo da resposta")
	}
	if tempF, ok := result["temp_F"]; !ok || tempF == 0 {
		t.Errorf("esperado temp_F no corpo da resposta")
	}
	if tempK, ok := result["temp_K"]; !ok || tempK == 0 {
		t.Errorf("esperado temp_K no corpo da resposta")
	}
}

func TestMultipleRequests(t *testing.T) {
	cep := "01001-000"
	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("Requisição #%d", i+1), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/weather?cep=%s", cep))
			if err != nil {
				t.Fatalf("erro na requisição: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("esperado 200, mas obteve %v", resp.StatusCode)
			}
		})
	}
}

func TestResponseTime(t *testing.T) {
	cep := "01001-000"
	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/weather?cep=%s", cep))
	if err != nil {
		t.Fatalf("erro na requisição: %v", err)
	}
	duration := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado 200, mas obteve %v", resp.StatusCode)
	}

	if duration > 2*time.Second {
		t.Errorf("tempo de resposta muito longo: %v", duration)
	}
}
