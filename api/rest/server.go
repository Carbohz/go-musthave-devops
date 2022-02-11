package rest

import (
	"context"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
	"net/http"
)

type APIServer struct {
	serverSvc  server.Processor
	httpServer *http.Server
}

func NewAPIServer(serverAddress string, serverSvc server.Processor) (*APIServer, error) {
	r := chi.NewRouter()
	setupRouters(r, serverSvc)

	srv := &APIServer{
		serverSvc: serverSvc,
		httpServer: &http.Server{
			Addr:    serverAddress,
			Handler: r,
		},
	}

	return srv, nil
}

func (s *APIServer) Run(ctx context.Context) error {

	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			log.Println("Done sub")
	//			return
	//		}
	//	}
	//}()

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %s", err)
	}

	return nil
}