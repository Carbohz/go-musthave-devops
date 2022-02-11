package rest

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/Carbohz/go-musthave-devops/service/server"
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
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
