package main

import "fmt"

type ID int

var (
	b bool
	c int
	d string
	e float64
	f ID = 1
)

func main() {
	fmt.Printf("o tipo do E é %T", e)
}
