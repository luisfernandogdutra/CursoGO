package main

import (
	"curso-go/matematica"
	"fmt"

	"github.com/google/uuid"
)

func main() {
	s := matematica.Soma(1, 2)
	fmt.Println("Resultado:", s)
	fmt.Println(matematica.A)
	//fmt.Println(matematica.b)

	carro := matematica.Carro{Marca: "Volks"}
	fmt.Println(carro)

	fmt.Println(carro.CarroAndando())

	fmt.Println(uuid.New())
}
