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
	"time"
)

const (
	htmlFile = "index.html"
	defaultAddress = "127.0.0.1:8080"
	defaultStoreInterval = 300 * time.Second // 300
	defaultStoreFile = "/tmp/devops-metrics-db.json" //"/tmp/devops-metrics-db.json"
	defaultRestore = true
)

//type Config struct {
//	Address string 				`env:"ADDRESS"`
//	StoreInterval time.Duration `env:"STORE_INTERVAL"`
//	StoreFile string 			`env:"STORE_FILE"`
//	Restore bool 				`env:"RESTORE"`
//}

func main() {
	// 1. Create config
	// 2. Prepare HTML page
	// 3. HandleInterrupts

	var cfg handler.Config

	cfg.StoreInterval = defaultStoreInterval
	cfg.StoreFile = defaultStoreFile
	cfg.Restore = true

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Address == "" {
		cfg.Address = defaultAddress
	}

	//if cfg.StoreInterval == 0 {
	//	cfg.StoreInterval = defaultStoreInterval
	//}
	//
	//if cfg.StoreFile == "" {
	//	cfg.StoreFile = defaultStoreFile
	//}

	//if !cfg.Restore {
	//	cfg.Restore = defaultRestore
	//}

	PrepareHTMLPage()
	//go RunServer(cfg)
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

	if cfg.Restore && cfg.StoreFile != "" {
		handler.LoadMetrics(cfg)
	}

	go RunServer(cfg)
	log.Println("awaiting signal")
	exitCode := <-exitChan
	log.Println("Saving metrics and exiting")
	handler.SaveMetrics(cfg)
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

func RunServer(cfg handler.Config) {
	//if cfg.Restore && cfg.StoreFile != "" {
	//	handler.LoadMetrics(cfg)
	//}

	if cfg.StoreInterval > 0 && cfg.StoreFile != "" {
		go handler.SaveMetrics(cfg)
	}

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

