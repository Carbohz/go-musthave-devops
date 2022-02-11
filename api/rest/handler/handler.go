package handler

import (
	"fmt"
	"github.com/markphelps/optional"
	"log"
	"net/http"
	"strconv"

	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
)

// не нужен -> router + mw + handler'ы
type Handler struct {
	serverSvc *server.Processor
	Router    *chi.Mux //отдельно хранить не нужно
}

func NewHandler(serverSvc *server.Processor) (*Handler, error) {
	r := chi.NewRouter()
	handler := &Handler{
		serverSvc: serverSvc,
		Router:    r,
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

	log.Printf("GaugeMetricHandler called. Requested metric %s with value %s", metricName, metricValue)

	gauge := model.Metric{Name: metricName, Type: model.KGauge, Value: optional.NewFloat64(value)}

	service := *h.serverSvc
	service.ProcessMetric(r.Context(), gauge)

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

	log.Printf("CounterMetricHandler called. Requested metric %s with value %s", metricName, metricValue)

	counter := model.Metric{Name: metricName, Type: model.KCounter, Delta: optional.NewInt64(value)}

	service := *h.serverSvc
	service.ProcessMetric(r.Context(), counter)

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

	log.Printf("SpecificMetricHandler called. Requested metric type is %s, name is %s", metricType, metricName)

	service := *h.serverSvc

	if m, found := service.GetMetric(metricName); found {
		w.WriteHeader(http.StatusOK)
		if delta, err := m.Delta.Get(); err == nil {
			w.Write([]byte(fmt.Sprint(delta)))
			log.Printf("Returned value from storage is %v", delta)
			return
		} else {
			value, _ := m.Value.Get()
			w.Write([]byte(fmt.Sprint(value)))
			log.Printf("Returned value from storage is %v", value)
			return
		}
	}
	log.Printf("No metric with type %s, name %s is storage", metricType, metricName)
	reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
	http.Error(w, reason, http.StatusNotFound)
}
