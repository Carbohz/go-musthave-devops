package main

import (
	"github.com/Carbohz/go-musthave-devops/internal/common"
	"github.com/Carbohz/go-musthave-devops/internal/handler"
	"github.com/Carbohz/go-musthave-devops/internal/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	htmlFile = "index.html"
)

func main() {
	cfg := server.CreateConfig()
	PrepareHTMLPage()
	db, err := handler.ConnectDB(cfg.DBPath)
	if db != nil {
		log.Printf("Connected to db: %s", cfg.DBPath)
	} else {
		log.Printf("DB connection error: %v", err)
	}
	exitChan := make(chan int, 1)
	go common.AwaitInterruptSignal(exitChan)
	go RunServer(cfg)
	exitCode := <-exitChan
	log.Println("Dumping metrics and exiting")
	handler.DumpMetricsImpl(cfg)
	if db != nil {
		db.Close()
	}
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

func RunServer(cfg server.Config) {
	// new function
	handler.PassSecretKey(cfg.Key)
	handler.PassServerConfig(cfg)

	if cfg.Restore && cfg.StoreFile != "" {
		handler.LoadMetrics(cfg)
	}

	if cfg.Restore && cfg.DBPath != "" {
		if err := handler.LoadStatsDB(); err != nil {
			log.Printf("Failed to load stats from Database: %v",err)
		}
	}

	//if cfg.Restore {
	//	if cfg.DBPath != "" {
	//		if err := handler.LoadStatsDB(); err != nil {
	//			log.Printf("Failed to load stats from Database: %v",err)
	//		}
	//	} else if cfg.StoreFile != "" {
	//		handler.LoadMetrics(cfg)
	//	}
	//}

	if cfg.StoreInterval > 0 && cfg.StoreFile != "" {
		go handler.DumpMetrics(cfg)
	}

	if cfg.DBPath != "" {
		if err := handler.InitDBTable(); err != nil {
			log.Printf("failed to init db tables: %v", err)
		}
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
