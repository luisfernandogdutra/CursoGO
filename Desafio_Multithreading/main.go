package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type Message struct {
	id  int64
	Msg string
}

func main() {
	c1 := make(chan Message)
	c2 := make(chan Message)
	var i int64 = 0
	// RabbitMQ
	go func() {
		for {
			atomic.AddInt64(&i, 1)
			msg := Message{i, BuscaViaCep("CEP VAI AQUI")}
			c1 <- msg
		}
	}()

	// Kafka
	go func() {
		for {
			atomic.AddInt64(&i, 1)
			msg := Message{i, BuscaBrasilApi("CEP VAI AQUI")}
			c2 <- msg
		}
	}()

	select {
	case msg := <-c1: // brasilapi
		fmt.Printf("Received from brasilapi: ID: %d - %s\n", msg.id, msg.Msg)

	case msg := <-c2: // viacep
		fmt.Printf("Received from viacep: ID: %d - %s\n", msg.id, msg.Msg)

	case <-time.After(time.Second * 1):
		println("timeout")
	}
}

func BuscaViaCep(cep string) string {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return "Erro ao consultar na viacep"
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Erro ao ler o retorno da viacep"
	}
	return string(body)
}

func BuscaBrasilApi(cep string) string {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		return "Erro ao consultar na brasilapi"
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Erro ao ler o retorno da brasilapi"
	}
	return string(body)
}
