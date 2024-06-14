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

func (cliente Cliente) Desativar() {
	cliente.Ativo = false
	fmt.Printf("O cliente %s foi desativado", cliente.Nome)
}

type Pessoa interface {
	Desativar()
}

type Empresa struct {
	Nome string
}

func (e Empresa) Desativar() {

}

func Desativacao(pessoa Pessoa) {
	pessoa.Desativar()
}

func main() {

	wesley := Cliente{
		Nome:  "Wesley",
		Idade: 30,
		Ativo: true,
	}

	minhaEmpresa := Empresa{}

	Desativacao(minhaEmpresa)
	Desativacao(wesley)
}
