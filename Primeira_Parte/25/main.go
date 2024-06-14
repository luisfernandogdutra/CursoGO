package main

import (
	"io"
	"net/http"
)

func main() {
	request, erro := http.Get("https://www.google.com")
	if erro != nil {
		panic(erro)
	}

	response, erro := io.ReadAll(request.Body)
	if erro != nil {
		panic(erro)
	}
	println(string(response))
	err := request.Body.Close()
	if err != nil {
		return
	}
}
