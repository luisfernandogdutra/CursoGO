package main

import "fmt"

type Endereco struct {
	Logradouro string
	Numero     int
	Cidade     string
	Estado     string
}

type Cliente struct {
	Nome  string
	Idade int
	Ativo bool
	Endereco
}

func (cliente Cliente) desativar() {
	cliente.Ativo = false
	fmt.Printf("O cliente %s foi desativado", cliente.Nome)
}

func main() {

	wesley := Cliente{
		Nome:  "Wesley",
		Idade: 30,
		Ativo: true,
	}

	wesley.desativar()
	// wesley.Endereco.Cidade = "Nome da cidade"

	// fmt.Printf("Nome: %s, Idade: %d, Ativo: %t", wesley.Nome, wesley.Idade, wesley.Ativo)
}
