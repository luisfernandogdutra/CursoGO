package main

import (
	"github.com/google/uuid"
	"github.com/luisfernandogdutra/CursoGO/Packaging/3/math"
)

func main() {
	m := math.NewMath()
	println(m.Add())
	println(uuid.New().String())
}
