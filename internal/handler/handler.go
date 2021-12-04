package handler

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"text/template"
)

var gaugeMetricsStorage = make(map[string]metrics.GaugeMetric)
var counterMetricsStorage = make(map[string]metrics.CounterMetric)
var HtmlTemplate *template.Template

func SetupRouters(r *chi.Mux) {
	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", GaugeMetricHandler)
		r.Post("/counter/{metricName}/{metricValue}", CounterMetricHandler)
		r.Post("/{metricName}/", NotFoundHandler)
		r.Post("/*", NotImplementedHandler)
	})
	r.Get("/value/{metricType}/{metricName}", SpecificMetricHandler)
	r.Get("/", AllMetricsHandler)
}

func GaugeMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "parsing error. Bad request", http.StatusBadRequest)
		return
	}
	gaugeMetricsStorage[metricName] = metrics.GaugeMetric{
		Base: metrics.Base{Name: metricName, Typename: metrics.Gauge},
		Value: value}
	w.WriteHeader(http.StatusOK)
}

func CounterMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	value, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		http.Error(w, "parsing error", http.StatusBadRequest)
		return
	}
	counterMetricsStorage[metricName] = metrics.CounterMetric{
		Base: metrics.Base{Name: metricName, Typename: metrics.Counter},
		Value: counterMetricsStorage[metricName].Value + value}
	w.WriteHeader(http.StatusOK)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unknown type", http.StatusNotImplemented)
}

func SpecificMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType == metrics.Counter {
		if val, found := counterMetricsStorage[metricName]; found {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(val.Value)))
			return
		}
		reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
		http.Error(w, reason, http.StatusNotFound)
		return
	}

	if metricType == metrics.Gauge {
		if val, found := gaugeMetricsStorage[metricName]; found {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(val.Value)))
			return
		}
		reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
		http.Error(w, reason, http.StatusNotFound)
	}
}

func AllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	renderData := map[string]interface{}{
		"gaugeMetrics": gaugeMetricsStorage,
		"counterMetrics": counterMetricsStorage,
	}
	HtmlTemplate.Execute(w, renderData)
}