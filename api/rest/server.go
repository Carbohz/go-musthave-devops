package rest

import (
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/Carbohz/go-musthave-devops/storage"
	"net/http"
)

var _ server.Processor = (*APIServer)(nil)

type APIServer struct {
	serverSvc server.Processor
	httpServer *http.Server
}

func NewAPIServer(serverAdress string, serverSvc server.Processor) (*APIServer, error) {
	return nil, nil
}