package rest

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

func GaugeMetricHandler(service server.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "parsing error. Bad request", http.StatusBadRequest)
			return
		}

		log.Printf("Requested to update storage for gauge metric %s to new value %s", metricName, metricValue)

		gauge := model.Metric{Name: metricName, Type: model.KGauge, Value: optional.NewFloat64(value)}

		service.ProcessMetric(r.Context(), gauge)

		w.WriteHeader(http.StatusOK)
	}
}

func CounterMetricHandler(service server.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "parsing error", http.StatusBadRequest)
			return
		}

		log.Printf("Requested to update storage for counter metric %s : need to add value %s to old one", metricName, metricValue)

		counter := model.Metric{Name: metricName, Type: model.KCounter, Delta: optional.NewInt64(value)}

		service.ProcessMetric(r.Context(), counter)

		w.WriteHeader(http.StatusOK)
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func UnknownTypeMetricHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unknown type", http.StatusNotImplemented)
}

func AllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func SpecificMetricHandler(service server.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		//log.Printf("SpecificMetricHandler called. Requested metric type is %s, name is %s", metricType, metricName)
		log.Printf("Requested to return metric %s of type %s from storage", metricName, metricType)

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
		log.Printf("No metric %s with type %s found in storage", metricName, metricType)
		reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
		http.Error(w, reason, http.StatusNotFound)
	}
}