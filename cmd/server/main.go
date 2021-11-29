package main

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/handler"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

const (
	host = "127.0.0.1"
	port = "8080"
)

func main() {
	RunServer()
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