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
	defaultStoreInterval = 2 * time.Second // 300
	defaultStoreFile = "tmp/devops-metrics-db.json" //"/tmp/devops-metrics-db.json"
	defaultRestore = true
)

type Config struct {
	Address string 				`env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile string 			`env:"STORE_FILE"`
	Restore bool 				`env:"RESTORE"`
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

	if cfg.StoreFile == "" {
		cfg.StoreFile = defaultStoreFile
	}

	if cfg.StoreInterval == 0 {
		cfg.StoreInterval = defaultStoreInterval
	}

	if !cfg.Restore {
		cfg.Restore = defaultRestore
	}

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
	go RunServer(cfg)
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
	if cfg.StoreInterval > 0 && cfg.StoreFile != "" {
		go metricsSaver(cfg)
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

func metricsSaver(cfg Config) {
	ticker := time.NewTicker(cfg.StoreInterval)
	for {
		<-ticker.C
		saveMetrics(cfg)
	}
}

func saveMetrics(cfg Config) {
	flags := os.O_WRONLY|os.O_CREATE|os.O_APPEND

	f, err := os.OpenFile(cfg.StoreFile, flags, 0777) //0644
	if err != nil {
		log.Fatal("cannot open file for writing: ", err)
	}
	defer f.Close()

	//if err := json.NewEncoder(f).Encode(statistics); err != nil {
	//	log.Fatal("cannot encode statistics: ", err)
	//}
	f.Write([]byte(`{"id":"llvm","type":"gauge","value":10}`))
}
