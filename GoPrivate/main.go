package main

import (
	"fmt"

	"github.com/luisfernandogdutra/fcutils-secret/pkg/events"
)

func main() {
	ed := events.NewEventDispatcher()
	fmt.Println(ed)
}
