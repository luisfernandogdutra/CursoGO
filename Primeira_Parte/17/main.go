package main

import "fmt"

type Cliente struct {
	nome string
}

func (c Cliente) andou() {
	c.nome = "Cliente Teste"
	fmt.Printf("O cliente %v andou \n", c.nome)
}

func main() {
	// cliente := Cliente{
	// 	nome: "Cliente",
	// }
	// cliente.andou()
	// fmt.Printf("O valor da struct com nome %v \n", cliente.nome)

	conta := Conta{saldo: 100}
	conta.simular(200)
	println(conta.saldo)
}

type Conta struct {
	saldo int
}

func NewConta() *Conta {
	return &Conta{saldo: 0}
}

func (c *Conta) simular(valor int) int {
	c.saldo += valor
	println(c.saldo)
	return c.saldo
}
