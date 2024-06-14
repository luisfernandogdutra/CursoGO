package main

import "fmt"

func main() {
	salarios := map[string]int{"Luis": 9000, "Joao": 5500, "Lucas": 2500}
	// fmt.Println(salarios["Luis"])
	// delete(salarios, "Luis")
	// salarios["Prolla"] = 35000

	// fmt.Println(salarios["Prolla"])

	// sal := make(map[string]int)
	// sal1 := map[string]int{}
	// sal1["Teste"] = 1000

	for nome, salario := range salarios {
		fmt.Printf("O salario de %s é %d\n", nome, salario)
	}

	for _, salario := range salarios {
		fmt.Printf("O salario é %d\n", salario)
	}
}
