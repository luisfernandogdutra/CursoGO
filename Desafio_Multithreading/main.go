package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type Message struct {
	id  int64
	Msg string
}

func main() {
	http.HandleFunc("/", BuscaCEPHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCEPHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c1 := make(chan Message)
	c2 := make(chan Message)
	var i int64 = 0
	var retorno string

	go func() {
		for {
			atomic.AddInt64(&i, 1)
			msg := Message{i, BuscaBrasilApi(cepParam)}
			c1 <- msg
		}
	}()

	go func() {
		for {
			atomic.AddInt64(&i, 1)
			msg := Message{i, BuscaViaCep(cepParam)}
			c2 <- msg
		}
	}()

	select {
	case msg := <-c1: // brasilapi
		retorno = "Received from brasilapi: ID: " + strconv.FormatInt(msg.id, 10) + " - " + msg.Msg
	case msg := <-c2: // viacep
		retorno = "Received from viacep: ID: " + strconv.FormatInt(msg.id, 10) + " - " + msg.Msg
	case <-time.After(time.Second * 1):
		retorno = "timeout"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(retorno)
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
