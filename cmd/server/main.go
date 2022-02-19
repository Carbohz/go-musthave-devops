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

	config, err := configsrv.NewCommonConfig()
	if err != nil {
		log.Fatalf("Failed to create common config: %v", err)
	}

	hybridStorageConfig := configsrv.NewHybridStorageConfig(config)
	storage, err := hybrid.NewMetricsStorage(hybridStorageConfig)
	if err != nil {
		log.Fatalf("Failed to create hybrid config: %v", err)
	}

	service := v1.NewService(storage)

	serverConfig := configsrv.NewServerConfig(config)
	apiServer, err := rest.NewAPIServer(serverConfig, service)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}
	defer apiServer.Close(ctx)

	go apiServer.Run(ctx)
	<-ctx.Done()
}
