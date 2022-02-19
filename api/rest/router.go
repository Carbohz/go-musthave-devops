package rest

import (
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func setupRouter(serverSvc server.Processor, key string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))

	// TODO! в chi есть mw, для проверки app-type json
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

	return r
}
