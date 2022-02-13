package main

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/api/rest"
	"github.com/Carbohz/go-musthave-devops/service/server"
	v1 "github.com/Carbohz/go-musthave-devops/service/server/v1"
	"github.com/Carbohz/go-musthave-devops/storage/inmemory"
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
	storage, _ := inmemory.NewMetricsStorage()
	// init server
	processor, _ := v1.NewService(storage) // serve
	// init apiServer
	apiServer, err := rest.NewAPIServer(config.Address, processor)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}

	//go func() {
	//	if err := apiServer.Run(ctx); err != nil {
	//		log.Fatalf("Failed to run a server: %v", err)
	//	}
	//}()
	//<-ctx.Done()

	go apiServer.Run(ctx)
	<-ctx.Done()
}
