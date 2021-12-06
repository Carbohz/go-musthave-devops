package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/Carbohz/go-musthave-devops/internal/handler"
	"github.com/go-chi/chi"
)

const (
	host     = "127.0.0.1"
	port     = "8080"
	htmlFile = "index.html"
)

func main() {
	PrepareHTMLPage()
	RunServer()
}

func PrepareHTMLPage() {
	page := strings.Join([]string{"cmd/server", htmlFile}, "/")
	bytes, err := os.ReadFile(page)
	if err != nil {
		log.Fatal("Error occurred while reading HTML file: ", err)
	}
	handler.HTMLTemplate, err = template.New("").Parse(string(bytes))
	if err != nil {
		log.Fatal("Error occurred while parsing HTML file: ", err)
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
