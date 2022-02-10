package rest

import (
	"context"
	"net/http"

	"github.com/Carbohz/go-musthave-devops/api/rest/handler"
	"github.com/Carbohz/go-musthave-devops/service/server"
)

//type Config struct {
//	Address       string        `env:"ADDRESS"`
//	StoreInterval time.Duration `env:"STORE_INTERVAL"`
//	StoreFile     string        `env:"STORE_FILE"`
//	Restore       bool          `env:"RESTORE"`
//	Key           string        `env:"KEY"`
//	DBPath        string        `env:"DATABASE_DSN"`
//}

type APIServer struct {
	serverSvc  server.Processor
	httpServer *http.Server
}

func NewAPIServer(serverAddress string, serverSvc server.Processor) (*APIServer, error) {
	//создаю Roter, регестрирую handler'ы
	//

	h, _ := handler.NewHandler(&serverSvc) // ?

	srv := &APIServer{
		serverSvc: serverSvc,
		httpServer: &http.Server{
			Addr:    serverAddress,
			Handler: h.Router, //должен создаваться в этой ф-ии
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
