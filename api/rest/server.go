package rest

import (
	"context"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	config server.Config
	serverSvc  server.Processor
	httpServer *http.Server
}

func NewAPIServer(config server.Config, serverSvc server.Processor) (*APIServer, error) {
	r := chi.NewRouter()
	setupRouters(r, serverSvc)

	srv := &APIServer{
		config: config,
		serverSvc: serverSvc,
		httpServer: &http.Server{
			Addr:    config.Address,
			Handler: r,
		},
	}

	return srv, nil
}

func (s *APIServer) Run(ctx context.Context) error {

	//go func() {
	//	storeTicker := time.NewTicker(s.config.StoreInterval)
	//	defer storeTicker.Stop()
	//	for {
	//		select {
	//		case <-storeTicker.C:
	//			s.serverSvc.Dump()
	//		//case <-ctx.Done():
	//		//	log.Println("Dumping and exiting")
	//		//	s.serverSvc.Dump()
	//		//	return
	//		}
	//	}
	//}()

	go func() {
		storeTicker := time.NewTicker(s.config.StoreInterval)
		defer storeTicker.Stop()
		for {
			<-storeTicker.C
			s.serverSvc.Dump()
		}
	}()

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %s", err)
	}

	return nil
}

func (s *APIServer) DumpBeforeExit() {
	log.Println("Dumping and exiting")
	s.serverSvc.Dump()
}