package main

import (
	"github.com/Carbohz/go-musthave-devops/internal/handler"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	htmlFile = "index.html"
	defaultAddress = "127.0.0.1:8080"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Address == "" {
		cfg.Address = defaultAddress
	}

	PrepareHTMLPage()
	RunServer(cfg)
}

func PrepareHTMLPage() {
	page := "cmd/server/" + htmlFile
	bytes, err := os.ReadFile(page)
	if err != nil {
		log.Fatal("Error occurred while reading HTML file: ", err)
	}
	handler.HTMLTemplate, err = template.New("").Parse(string(bytes))
	if err != nil {
		log.Fatal("Error occurred while parsing HTML file: ", err)
	}
}

func RunServer(cfg Config) {
	r := chi.NewRouter()
	handler.SetupRouters(r)
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}
	server.SetKeepAlivesEnabled(false)
	log.Printf("Listening on address %s", cfg.Address)
	log.Fatal(server.ListenAndServe())
}
