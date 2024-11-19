package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Estrutura para armazenar o relatório do teste de carga
type LoadTestReport struct {
	TotalRequests      int
	SuccessfulRequests int
	StatusCodeCounts   map[int]int
	TotalTime          time.Duration
}

func main() {
	// Definição dos parâmetros via CLI
	url := flag.String("url", "", "URL do serviço a ser testado")
	totalRequests := flag.Int("requests", 1000, "Número total de requests")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	flag.Parse()

	// Verificar se a URL foi fornecida
	if *url == "" {
		log.Fatal("A URL do serviço é obrigatória")
	}

	// Iniciar o teste de carga
	report := runLoadTest(*url, *totalRequests, *concurrency)

	// Gerar o relatório
	printReport(report)
}

// Função para realizar o teste de carga
func runLoadTest(url string, totalRequests, concurrency int) LoadTestReport {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var report LoadTestReport
	report.StatusCodeCounts = make(map[int]int)

	startTime := time.Now()

	// Função para realizar uma única requisição
	doRequest := func() {
		defer wg.Done()

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Erro ao fazer request: %v", err)
			return
		}
		defer resp.Body.Close()

		// Bloqueio para garantir thread-safety ao atualizar o relatório
		mu.Lock()
		report.TotalRequests++
		report.StatusCodeCounts[resp.StatusCode]++
		if resp.StatusCode == http.StatusOK {
			report.SuccessfulRequests++
		}
		mu.Unlock()
	}

	// Lançar as requisições com concorrência
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go doRequest()
		if i%concurrency == 0 {
			wg.Wait() // Espera um lote de requisições terminar antes de continuar
		}
	}

	wg.Wait() // Espera todos os requests terminarem
	report.TotalTime = time.Since(startTime)

	return report
}

// Função para imprimir o relatório após o teste
func printReport(report LoadTestReport) {
	fmt.Println("Relatório do Teste de Carga:")
	fmt.Printf("Tempo total de execução: %v\n", report.TotalTime)
	fmt.Printf("Total de Requests realizados: %d\n", report.TotalRequests)
	fmt.Printf("Total de Requests com sucesso (HTTP 200): %d\n", report.SuccessfulRequests)
	fmt.Println("Distribuição de Códigos de Status HTTP:")
	for code, count := range report.StatusCodeCounts {
		fmt.Printf("  %d: %d\n", code, count)
	}
}
