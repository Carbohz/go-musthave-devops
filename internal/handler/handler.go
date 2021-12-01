package handler

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type gauge struct {
	v float64
}

type counter struct {
	v int64
}

var gaugeMetricsStorage = make(map[string]gauge)
var counterMetricsStorage = make(map[string]counter)

func SetupRouters(r *chi.Mux) {
	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", GaugeMetricHandler)
		r.Post("/counter/{metricName}/{metricValue}", CounterMetricHandler)
		r.Post("/{metricName}/", NotFoundHandler)
		r.Post("/*", NotImplementedHandler)
	})
	r.Get("/value/*", SpecificMetricHandler)
	r.Get("/", AllMetricsHandler)
}

func GaugeMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricName, metricValue := GetRequestBody(r)
	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "parsing error. Bad request", http.StatusBadRequest)
		return
	}
	gaugeMetricsStorage[metricName] = gauge{v: value}
	w.WriteHeader(http.StatusOK)
}

func CounterMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricName, metricValue := GetRequestBody(r)
	value, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		http.Error(w, "parsing error", http.StatusBadRequest)
		return
	}
	counterMetricsStorage[metricName] = counter{counterMetricsStorage[metricName].v + value}
	w.WriteHeader(http.StatusOK)
}

func GetRequestBody(r *http.Request) (string, string) {
	uri := r.RequestURI
	tokens := strings.Split(uri, "/")
	metricName := tokens[3]
	metricValue := tokens[4]
	return metricName, metricValue
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unknown type", http.StatusNotImplemented)
}

func SpecificMetricHandler(w http.ResponseWriter, r *http.Request) {
	tokens := strings.Split(r.URL.Path, "/")
	metricType := tokens[2]
	metricName := tokens[3]

	if metricType == "counter" {
		if val, found := counterMetricsStorage[metricName]; found {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(val.v)))
		} else {
			reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
			http.Error(w, reason, http.StatusNotFound)
		}
	}

	if metricType == "gauge" {
		if val, found := gaugeMetricsStorage[metricName]; found {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(val.v)))
		} else {
			reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
			http.Error(w, reason, http.StatusNotFound)
		}
	}
}

func AllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// htmlFile := "index.html"
	htmlFile := "D:\\Go\\yandex-praktikum\\Sprint1\\net_http\\increment1\\go-musthave-devops2\\cmd\\server\\index.html"
	htmlPage, err := os.ReadFile(htmlFile)
	if err != nil {
		log.Println("File reading error:", err)
		os.Exit(1)
	}

	renderData := map[string]interface{}{
		"gaugeMetrics": gaugeMetricsStorage,
		"counterMetrics": counterMetricsStorage,
	}
	tmpl := template.Must(template.New("").Parse(string(htmlPage)))
	tmpl.Execute(w, renderData)
}
