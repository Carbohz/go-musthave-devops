package rest

import (
	"context"
	"fmt"
	configsrv "github.com/Carbohz/go-musthave-devops/config/server"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	config configsrv.ServerConfig
	serverSvc  server.Processor
	httpServer *http.Server
}

func NewAPIServer(config configsrv.ServerConfig, serverSvc server.Processor) (*APIServer, error) {
	r := chi.NewRouter()
	// д.б. в setupRouters
	r.Use(middleware.Compress(5))

	setupRouters(r, serverSvc, config.Key)

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
	// TODO! не использую ctx
	// TODO! goroutine с завершением ctx

	go func() {
		storeTicker := time.NewTicker(s.config.StoreInterval)
		defer storeTicker.Stop()
		for {
			<-storeTicker.C
			s.serverSvc.Dump(ctx)
		}
	}()

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %s", err)
	}

	return nil
}

// DumpBeforeExit() -> defer Close()
func (s *APIServer) DumpBeforeExit() {
	// TODO! Здесь можно выключить, тогда в Run не нужен ctx
	log.Println("Dumping and exiting")
	ctx := context.Background()
	s.serverSvc.Dump(ctx)
}