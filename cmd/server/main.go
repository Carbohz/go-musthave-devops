package main

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/handler"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	host = "127.0.0.1"
	port = "8080"
	htmlFile = "index.html"
)

func main() {
	PrepareHtmlFile()
	RunServer()
}

func PrepareHtmlFile() {
	bytes, err := os.ReadFile(htmlFile)
	if err != nil {
		log.Fatal(err)
	}
	handler.HtmlTemplate, err = template.New("").Parse(string(bytes))
	if err != nil {
		log.Fatal(err)
	}
}

func RunServer() {
	r := chi.NewRouter()
	handler.SetupRouters(r)
	addr := fmt.Sprintf("%s:%s", host, port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	server.SetKeepAlivesEnabled(false)
	log.Printf("Listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}
