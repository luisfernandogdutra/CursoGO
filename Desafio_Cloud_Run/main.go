package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type CepResponse struct {
	Cidade string `json:"localidade"`
}

func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{5}-\d{3}$`, cep)
	return match
}

func getCityFromCEP(cep string) (string, error) {
	cep = regexp.MustCompile(`[^0-9]`).ReplaceAllString(cep, "")
	if len(cep) != 8 {
		return "", errors.New("CEP inválido")
	}

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar o CEP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("não foi possível encontrar o CEP")
	}

	var cepData CepResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler a resposta do CEP: %v", err)
	}

	if err := json.Unmarshal(body, &cepData); err != nil {
		return "", fmt.Errorf("erro ao parsear o JSON do CEP: %v", err)
	}

	if cepData.Cidade == "" {
		return "", errors.New("cidade não encontrada para o CEP")
	}

	return cepData.Cidade, nil
}

func getWeather(city string) (float64, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=bad4fe6d4148402daa712525242412&q=%s", encodedCity)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("erro ao consultar o clima: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("não foi possível obter o clima")
	}

	var weatherData WeatherResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("erro ao ler a resposta do clima: %v", err)
	}

	if err := json.Unmarshal(body, &weatherData); err != nil {
		return 0, fmt.Errorf("erro ao parsear o JSON do clima: %v", err)
	}

	return weatherData.Current.TempC, nil
}

func convertTemperatures(tempC float64) (float64, float64, float64) {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273
	return tempC, tempF, tempK
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !isValidCEP(cep) {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	city, err := getCityFromCEP(cep)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message": "%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	tempC, err := getWeather(city)
	if err != nil {
		http.Error(w, `{"message": "can not find zipcode"}`, http.StatusNotFound)
		return
	}

	tempC, tempF, tempK := convertTemperatures(tempC)

	response := map[string]float64{
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/weather", handleRequest)
	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
