package main

import (
	"net/http"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.Handle("/blog", blog{title: "My Blog"})
	http.ListenAndServe(":8080", mux)
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Hello, World!"))
}

type blog struct {
	title string
}

func (b blog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Blog"))
}
