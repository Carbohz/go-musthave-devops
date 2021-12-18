package main

import (
	"flag"
	"github.com/Carbohz/go-musthave-devops/internal/handler"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
)

const (
	htmlFile 				= "index.html"
	defaultAddress 			= "127.0.0.1:8080"
	defaultStoreInterval 	= 300 * time.Second
	defaultStoreFile 		= "/tmp/devops-metrics-db.json"
	defaultRestore 			= true
)

func main() {
	cfg := CreateConfig()
	PrepareHTMLPage()
	exitChan := make(chan int, 1)
	go AwaitInterruptSignal(exitChan)
	go RunServer(cfg)
	exitCode := <-exitChan
	log.Println("Dumping metrics and exiting")
	handler.DumpMetricsImpl(cfg)
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
	if cfg.Restore && cfg.StoreFile != "" {
		handler.LoadMetrics(cfg)
	}

	if cfg.StoreInterval > 0 && cfg.StoreFile != "" {
		go handler.DumpMetrics(cfg)
	}

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	handler.SetupRouters(r)
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}
	server.SetKeepAlivesEnabled(false)
	log.Printf("Listening on address %s", cfg.Address)
	log.Fatal(server.ListenAndServe())
}

func CreateConfig() handler.Config {
	var cfg handler.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server is running with environment variables: %+v", cfg)

	addressFlagPtr := flag.String("a", defaultAddress, "set address of server")
	storeIntervalFlagPtr := flag.Duration("i", defaultStoreInterval, "set server's metrics store interval")
	storeFileFlagPtr := flag.String("f", defaultStoreFile, "set file where metrics are stored")
	restoreFlagPtr := flag.Bool("r", defaultRestore, "choose whether to restore server metrics from file")
	flag.Parse()
	log.Printf("Server is running with command line flags: Address %v, Store Interval %v, Store File %v, Restore %v",
		*addressFlagPtr, *storeIntervalFlagPtr, *storeFileFlagPtr, *restoreFlagPtr)

	_, isSet := os.LookupEnv("ADDRESS")
	if !isSet {
		cfg.Address = *addressFlagPtr
	}

	_, isSet = os.LookupEnv("STORE_INTERVAL")
	if !isSet {
		cfg.StoreInterval = *storeIntervalFlagPtr
	}

	_, isSet = os.LookupEnv("STORE_FILE")
	if !isSet {
		cfg.StoreFile = *storeFileFlagPtr
	}

	_, isSet = os.LookupEnv("RESTORE")
	if !isSet {
		cfg.Restore = *restoreFlagPtr
	}

	return cfg
}

func AwaitInterruptSignal(exitChan chan<- int) {
	log.Println("Awaiting interrupt signal")

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		s := <-signalChanel
		switch s {
		case syscall.SIGINT:
			log.Printf("%s SIGINT signal triggered.", s)
			exitChan <- 1

		case syscall.SIGTERM:
			log.Printf("%s SIGTERM signal triggered.", s)
			exitChan <- 2

		case syscall.SIGQUIT:
			log.Printf("%s SIGQUIT signal triggered.", s)
			exitChan <- 3

		default:
			log.Printf("%s UNKNOWN signal triggered.", s)
			exitChan <- 1
		}
	}()
}