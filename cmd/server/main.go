package main

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/api/rest"
	v1 "github.com/Carbohz/go-musthave-devops/service/server/v1"
	"github.com/Carbohz/go-musthave-devops/storage/inmemory"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	// init storage
	storage, _ := inmemory.NewMetricsStorage()
	// init server
	processor, _ := v1.NewService(storage) // serve
	// init apiServer
	apiServer, err := rest.NewAPIServer("127.0.0.1:8080", processor)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}

	if err := apiServer.Run(ctx); err != nil {
		log.Fatalf("Failed to run a server: %v", err)
	}
}
