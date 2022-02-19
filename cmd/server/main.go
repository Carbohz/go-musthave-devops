package main

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/api/rest"
	configsrv "github.com/Carbohz/go-musthave-devops/config/server"
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

	config, err := configsrv.NewCommonConfig()
	if err != nil {
		log.Fatalf("Failed to create common config: %v", err)
	}

	// hybrid storage
	hybridStorageConfig := configsrv.NewHybridStorageConfig(config)
	storage, err := hybrid.NewMetricsStorage(hybridStorageConfig)
	if err != nil {
		log.Println("Failed to create hybrid config")
	}

	processor, _ := v1.NewService(storage)
	// _ -> err

	serverConfig := configsrv.NewServerConfig(config)
	apiServer, err := rest.NewAPIServer(serverConfig, processor)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}
	defer apiServer.DumpBeforeExit() // ctx

	go apiServer.Run(ctx)
	<-ctx.Done()
}
