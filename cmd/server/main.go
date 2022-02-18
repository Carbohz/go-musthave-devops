package main

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/api/rest"
	"github.com/Carbohz/go-musthave-devops/service/server"
	v1 "github.com/Carbohz/go-musthave-devops/service/server/v1"
	"github.com/Carbohz/go-musthave-devops/storage/hybrid"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, ctxCancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer ctxCancel()
	// нужен timeout

	config := server.CreateConfig()

	// hybrid storage
	hybridConfig := hybrid.Config{
		StoreInterval: config.StoreInterval,
		StoreFile: config.StoreFile,
		Restore: config.Restore,
		DBPath: config.DBPath}
	storage, err := hybrid.NewMetricsStorage(hybridConfig)
	if err != nil {
		log.Println("Failed to create hybrid config")
	}

	processor, _ := v1.NewService(storage)
	// _ -> err

	apiServer, err := rest.NewAPIServer(config, processor)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}
	defer apiServer.DumpBeforeExit() // ctx

	go apiServer.Run(ctx)
	<-ctx.Done()
}
