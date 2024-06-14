package matematica

func Soma[T int | float64](a, b T) T {
	return a + b
}

var A int = 10

// var b int = 5

type Carro struct {
	Marca string
}

// type carro struct {
// 	Marca string
// }

func (c Carro) CarroAndando() string {
	return "Carro Andando"
}
