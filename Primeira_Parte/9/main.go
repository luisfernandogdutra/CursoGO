package main

import (
	"fmt"
)

func main() {
	fmt.Println(sum(1, 3, 45, 12, 145, 10, 213, 534, 123, 1241, 123, 22))
}

func sum(numeros ...int) int {
	total := 0
	for _, numero := range numeros {
		total += numero
	}

	return total
}
