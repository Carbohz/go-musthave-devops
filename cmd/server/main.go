package main

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/api/rest"
	"github.com/Carbohz/go-musthave-devops/service/server"
	v1 "github.com/Carbohz/go-musthave-devops/service/server/v1"
	"github.com/Carbohz/go-musthave-devops/storage/filebased"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, ctxCancel  := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer ctxCancel()

	config := server.CreateConfig()

	// init storage
	//storage, _ := inmemory.NewMetricsStorage()
	storageConfig := filebased.Config{
		StoreInterval: config.StoreInterval,
		StoreFile: config.StoreFile,
		Restore: config.Restore}
	storage, _ := filebased.NewMetricsStorage(storageConfig)
	// init server
	processor, _ := v1.NewService(storage) // serve
	// init apiServer
	apiServer, err := rest.NewAPIServer(config, processor)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}

	go apiServer.Run(ctx)
	<-ctx.Done()
	log.Println("Dumping and exiting")
	processor.Dump()
}
