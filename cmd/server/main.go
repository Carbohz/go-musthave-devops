package main

import (
	"github.com/Carbohz/go-musthave-devops/internal/handler"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	go RunServer(cfg)
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT, // interrupt
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int, 1)
	go func() {
		s := <-signalChanel
		switch s {
		case syscall.SIGINT:
			log.Printf("%s signal triggered.", s)
			exitChan <- 1

		case syscall.SIGTERM:
			log.Printf("%s signal triggered.", s)
			exitChan <- 2

		case syscall.SIGQUIT:
			log.Printf("%s signal triggered.", s)
			exitChan <- 3

		default:
			log.Printf("%s signal triggered.", s)
			exitChan <- 1
		}
	}()
	//RunServer(cfg)
	//exitCode := <-exitChan
	//os.Exit(exitCode)
	log.Println("awaiting signal")
	exitCode := <-exitChan
	log.Println("exiting")
	os.Exit(exitCode)
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
