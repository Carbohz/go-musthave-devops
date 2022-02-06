package handler

import (
	"context"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	serverSvc *server.Processor
	Router *chi.Mux
}

func NewHandler(serverSvc *server.Processor) (*Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handler := &Handler{
		serverSvc: serverSvc,
		Router: r,
	}

	handler.setupRouters()

	log.Println("Created NewHandler")
	return handler, nil
}

func (h *Handler) setupRouters() {
	r := h.Router
	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", h.GaugeMetricHandler)
		r.Post("/counter/{metricName}/{metricValue}", h.CounterMetricHandler)
	})
	r.Get("/", h.AllMetricsHandler)

	log.Println("Routers set up")
}


func (h *Handler) GaugeMetricHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GaugeMetricHandler called")
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

	service := *h.serverSvc
	service.ProcessGaugeMetric(ctx, gauge)
	//h.serverSvc.ProcessGaugeMetric(ctx, gauge)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CounterMetricHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CounterMetricHandler called")
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

	service := *h.serverSvc
	service.ProcessCounterMetric(ctx, counter)
	//h.serverSvc.ProcessCounterMetric(ctx, counter)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AllMetricsHandler called")

	//w.Header().Set("Content-Type", "text/html")
	//renderData := map[string]interface{}{
	//	"gaugeMetrics":   gaugeMetricsStorage,
	//	"counterMetrics": counterMetricsStorage,
	//}
	//HTMLTemplate.Execute(w, renderData)

	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}