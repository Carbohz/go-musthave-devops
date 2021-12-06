package main

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/handler"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const (
	host = "127.0.0.1"
	port = "8080"
	htmlFile = "index.html"
)

func main() {
	PrepareHTMLPage()
	RunServer()
}

func PrepareHTMLPage() {
	path, _ := os.Getwd()
	abs := strings.Join([]string{path, "cmd/server", htmlFile}, "\\")
	fmt.Println(abs)
	bytes, err := os.ReadFile(abs) // htmlFile
	if err != nil {
		fmt.Println("Error occurred while reading HTML file: ", err)
		log.Fatal("Error occurred while reading HTML file: ", err)
	}
	handler.HTMLTemplate, err = template.New("").Parse(string(bytes))
	if err != nil {
		fmt.Println("Error occurred while parsing HTML file: ", err)
		log.Fatal("Error occurred while parsing HTML file: ", err)
	}
}

func RunServer() {
	log.Println("Running server")
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
