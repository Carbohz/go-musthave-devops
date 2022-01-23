package main

import (
	"github.com/Carbohz/go-musthave-devops/internal/common"
	"github.com/Carbohz/go-musthave-devops/internal/server"
	"log"
	"os"
	"text/template"
)

const (
	htmlFile = "index.html"
)

func main() {
	cfg := server.CreateConfig()
	PrepareHTMLPage()
	instance := server.CreateInstance(cfg)
	exitChan := make(chan int, 1)
	go common.AwaitInterruptSignal(exitChan)
	//go RunServer(cfg)
	go instance.RunInstance()
	defer instance.BeforeShutDown()
	exitCode := <-exitChan
	//log.Println("Dumping metrics and exiting")
	//handler.DumpMetricsImpl(cfg)

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

//func RunServer(cfg server.Config) {
//	if cfg.Restore && cfg.StoreFile != "" {
//		handler.LoadMetrics(cfg)
//	}
//
//	if cfg.StoreInterval > 0 && cfg.StoreFile != "" {
//		go handler.DumpMetrics(cfg)
//	}
//
//	// new function
//	handler.PassSecretKey(cfg.Key)
//
//	//db, err := handler.ConnectDB(cfg.DBPath)
//
//	r := chi.NewRouter()
//	r.Use(middleware.Compress(5))
//	handler.SetupRouters(r)
//	server := &http.Server{
//		Addr:    cfg.Address,
//		Handler: r,
//	}
//	server.SetKeepAlivesEnabled(false)
//	log.Printf("Listening on address %s", cfg.Address)
//	log.Fatal(server.ListenAndServe())
//}
