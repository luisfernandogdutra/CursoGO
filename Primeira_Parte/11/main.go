package main

import "fmt"

type Cliente struct {
	Nome  string
	Idade int
	Ativo bool
}

func main() {

	wesley := Cliente{
		Nome:  "Wesley",
		Idade: 30,
		Ativo: true,
	}

	wesley.Ativo = false

	fmt.Printf("Nome: %s, Idade: %d, Ativo: %t", wesley.Nome, wesley.Idade, wesley.Ativo)
}
