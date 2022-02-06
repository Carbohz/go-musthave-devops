package handler

import (
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
)

type Handler struct {
	Server *server.Server
	Router *chi.Mux
}
