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
		// в {} можно добавить regex
		r.Post("/gauge/{metricName}/{metricValue}", GaugeMetricHandler(serverSvc))
		r.Post("/gauge/*", InvalidNameAndValueHandler)
		r.Post("/counter/{metricName}/{metricValue}", CounterMetricHandler(serverSvc))
		r.Post("/counter/*", InvalidNameAndValueHandler)
		r.Post("/{param}/*", UnknownMetricTypeHandler) // ?
		r.Post("/", UpdateMetricsJSONHandler(serverSvc, key))
	})
	r.Post("/updates/", UpdatesMetricsJSONHandler(serverSvc, key))
	// value можно сгруппировать
	r.Post("/value/", GetMetricsJSONHandler(serverSvc, key))
	r.Get("/value/{metricType}/{metricName}", SpecificMetricHandler(serverSvc))
}
