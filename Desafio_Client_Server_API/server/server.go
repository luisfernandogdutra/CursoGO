package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FinancialQuote struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type Cotation struct {
	ID    int `gorm:"primaryKey"`
	Value string
	Date  string
}

type Dolar struct {
	Dolar string `json:"Dólar"`
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")
	w.Write([]byte(requestCotacao()))
}

func requestCotacao() string {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequest("GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL/", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return "Error: " + err.Error()
	}

	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}

	log.Println("Resposta do get -> " + string(res))
	var data FinancialQuote
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}

	log.Println("Vai salvar no banco de dados")
	salvaCotacao(data)

	u, err := json.Marshal(Dolar{
		Dolar: data.Usdbrl.Bid,
	})
	if err != nil {
		panic(err)
	}

	log.Println("Retorno: " + string(u))
	return string(u)
}

func salvaCotacao(financial FinancialQuote) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Println("Erro ao conectar no banco de dados")
		panic(err)
	}
	err = db.AutoMigrate(&Cotation{})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	cotation := &Cotation{
		Value: financial.Usdbrl.Bid,
		Date:  financial.Usdbrl.CreateDate,
	}

	db.Session(&gorm.Session{SkipDefaultTransaction: true, Context: ctx}).Create(cotation)
}
