package rest

import (
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
)

func setupRouters(r *chi.Mux, serverSvc server.Processor, key string) {
	// здесь можно создать router с возвращением его

	// в chi есть mw, для проверки app-type json
	r.Get("/", AllMetricsHandler)
	r.Get("/ping", PingDBHandler(serverSvc))

	r.Route("/update", func(r chi.Router) {
		r.Post("/", UpdateMetricsJSONHandler(serverSvc, key))

		r.Route("/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
			r.Use(metricTypeValidator, metricNameValidator, metricValueValidator)
			r.Post("/", URLMetricHandler(serverSvc))
		})
	})
	r.Post("/updates/", UpdatesMetricsJSONHandler(serverSvc, key))

	r.Route("/value", func(r chi.Router) {
		r.Post("/", GetMetricsJSONHandler(serverSvc, key))

		r.Route("/{metricType}/{metricName}", func(r chi.Router) {
			r.Use(metricTypeValidator, metricNameValidator)
			r.Get("/", SpecificMetricHandler(serverSvc))
		})
	})
}
