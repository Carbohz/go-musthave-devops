package handler

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	serverSvc server.Processor
	Router *chi.Mux
}

func NewHandler(serverSvc server.Processor) (*Handler, error) {
	r := chi.NewRouter()

	handler := &Handler{
		serverSvc: serverSvc,
		Router: r,
	}

	setupRouters(r, handler)

	log.Println("Created NewHandler")
	return handler, nil
}

func setupRouters(r *chi.Mux, h *Handler) {
	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", h.GaugeMetricHandler)
		r.Post("/counter/{metricName}/{metricValue}", h.CounterMetricHandler)
		//r.Post("/{metricName}/", NotFoundHandler)
		//r.Post("/*", UnknownTypeMetricHandler)
		//r.Post("/", UpdateMetricsJSONHandler)
	})
	//r.Post("/value/", GetMetricsJSONHandler)
	//r.Get("/value/{metricType}/{metricName}", SpecificMetricHandler)
	//r.Get("/", AllMetricsHandler)
	//r.Get("/ping", PingDBHandler)
	//r.Post("/updates/", UpdatesMetricsJSONHandler)

	log.Println("Routers set up")
}


func (h *Handler) GaugeMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "parsing error. Bad request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	request := models.GaugeMetricRequest{MType: model.Gauge, Name: metricName, Value: value}
	gauge := request.ToModelGaugeMetric()

	h.serverSvc.ProcessGaugeMetric(ctx, gauge)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CounterMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	value, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		http.Error(w, "parsing error", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	request := models.CounterMetricRequest{MType: model.Gauge, Name: metricName, Value: value}
	counter := request.ToModelCounterMetric()

	h.serverSvc.ProcessCounterMetric(ctx, counter)

	w.WriteHeader(http.StatusOK)
}