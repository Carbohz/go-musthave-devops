package rest

import "github.com/go-chi/chi"

func SetupRouters(r *chi.Mux) {
	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", GaugeMetricHandler)
		r.Post("/counter/{metricName}/{metricValue}", CounterMetricHandler)
		//r.Post("/{metricName}/", NotFoundHandler)
		//r.Post("/*", UnknownTypeMetricHandler)
		//r.Post("/", UpdateMetricsJSONHandler)
	})
	//r.Post("/value/", GetMetricsJSONHandler)
	//r.Get("/value/{metricType}/{metricName}", SpecificMetricHandler)
	//r.Get("/", AllMetricsHandler)
	//r.Get("/ping", PingDBHandler)
	//r.Post("/updates/", UpdatesMetricsJSONHandler)
}

