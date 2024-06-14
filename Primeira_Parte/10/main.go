package main

import (
	"fmt"
)

func main() {

	total := func() int {
		return sum(1, 3, 45, 12, 145, 10, 213, 534, 123, 1241, 123, 22) * 2
	}()

	fmt.Println(total)
}

func sum(numeros ...int) int {
	total := 0
	for _, numero := range numeros {
		total += numero
	}

	return total
}
