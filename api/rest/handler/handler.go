package handler

import (
	"context"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

type Handler struct {
	serverSvc *server.Processor
	Router *chi.Mux
}

func NewHandler(serverSvc *server.Processor) (*Handler, error) {
	r := chi.NewRouter()
	handler := &Handler{
		serverSvc: serverSvc,
		Router: r,
	}
	handler.setupRouters()
	return handler, nil
}

func (h *Handler) setupRouters() {
	r := h.Router
	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", h.GaugeMetricHandler)
		r.Post("/counter/{metricName}/{metricValue}", h.CounterMetricHandler)
		r.Post("/{metricName}/", h.NotFoundHandler)
		r.Post("/*", h.UnknownTypeMetricHandler)
	})
	r.Get("/value/{metricType}/{metricName}", h.SpecificMetricHandler)
	r.Get("/", h.AllMetricsHandler)
}


func (h *Handler) GaugeMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "parsing error. Bad request", http.StatusBadRequest)
		return
	}

	request := models.GaugeMetricRequest{MType: model.Gauge, Name: metricName, Value: value}
	gauge := request.ToModelGaugeMetric()

	ctx := context.Background()
	service := *h.serverSvc
	service.ProcessGaugeMetric(ctx, gauge)
	//h.serverSvc.ProcessGaugeMetric(ctx, gauge)

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

	request := models.CounterMetricRequest{MType: model.Gauge, Name: metricName, Value: value}
	counter := request.ToModelCounterMetric()

	ctx := context.Background()
	service := *h.serverSvc
	service.ProcessCounterMetric(ctx, counter)
	//h.serverSvc.ProcessCounterMetric(ctx, counter)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func (h *Handler) UnknownTypeMetricHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unknown type", http.StatusNotImplemented)
}

func (h *Handler) AllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func (h *Handler) SpecificMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	service := *h.serverSvc

	if metricType == model.Counter {
		if value, found := service.GetCounterMetric(metricName); found {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(value)))
			return
		}
		reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
		http.Error(w, reason, http.StatusNotFound)
		return
	}

	if metricType == model.Gauge {
		if value, found := service.GetGaugeMetric(metricName); found {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(value)))
			return
		}
		reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
		http.Error(w, reason, http.StatusNotFound)
	}
}