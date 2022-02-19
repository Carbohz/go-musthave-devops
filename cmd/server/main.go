package main

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/api/rest"
	configsrv "github.com/Carbohz/go-musthave-devops/config/server"
	v1 "github.com/Carbohz/go-musthave-devops/service/server/v1"
	"github.com/Carbohz/go-musthave-devops/storage/filebased"
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

	//// hybrid storage
	//hybridStorageConfig := configsrv.NewHybridStorageConfig(config)
	//storage, err := hybrid.NewMetricsStorage(hybridStorageConfig)
	//if err != nil {
	//	log.Fatalf("Failed to create hybrid config: %v", err)
	//	//log.Printf("Failed to create hybrid config: %v", err)
	//}

	// fileBased storage
	fileBasedConfig := configsrv.NewFileBasedStorageConfig(config)
	storage, err := filebased.NewMetricsStorage(fileBasedConfig)
	if err != nil {
		log.Fatalf("Failed to create filebased config: %v", err)
		//log.Printf("Failed to create filebased config: %v", err)
	}

	service := v1.NewService(storage)

	serverConfig := configsrv.NewServerConfig(config)
	apiServer, err := rest.NewAPIServer(serverConfig, service)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}
	defer apiServer.DumpBeforeExit(ctx) // ctx

	go apiServer.Run(ctx)
	<-ctx.Done()
}
