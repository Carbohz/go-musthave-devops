package rest

import (
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
)

func setupRouters(r *chi.Mux, serverSvc  server.Processor) {
	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", GaugeMetricHandler(serverSvc))
		r.Post("/counter/{metricName}/{metricValue}", CounterMetricHandler(serverSvc))
		r.Post("/{metricName}/", NotFoundHandler)
		r.Post("/*", UnknownTypeMetricHandler)
		//r.Post("/", UpdateMetricsJSONHandler(serverSvc))
	})
	//r.Post("/value/", GetMetricsJSONHandler(serverSvc))
	r.Get("/value/{metricType}/{metricName}", SpecificMetricHandler(serverSvc))
	r.Get("/", AllMetricsHandler)
}