package main

func SomaInteiro(m map[string]int) int {
	var soma int
	for _, v := range m {
		soma += v
	}

	return soma
}

func SomaFloat(m map[string]float64) float64 {
	var soma float64
	for _, v := range m {
		soma += v
	}

	return soma
}

type MyNumber int

type Number interface {
	~int | ~float64
}

func Soma[T Number](m map[string]T) T {
	var soma T
	for _, v := range m {
		soma += v
	}

	return soma
}

func Compara[T comparable](a T, b T) bool {
	if a == b {
		return true
	}

	return false
}

func main() {
	m := map[string]int{"Luis": 9000, "João": 2000, "Lucas": 1500}
	m2 := map[string]float64{"Luis": 9000.10, "João": 2000.20, "Lucas": 1500.05}
	// println(SomaInteiro(m))
	// println(SomaFloat(m2))

	m3 := map[string]MyNumber{"Luis": 9000, "João": 2000, "Lucas": 1500}

	println(Soma(m))
	println(Soma(m2))
	println(Soma(m3))

	println(Compara(10, 10.0))
}
