package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"github.com/go-chi/chi"
)

var gaugeMetricsStorage = make(map[string]metrics.GaugeMetric)
var counterMetricsStorage = make(map[string]metrics.CounterMetric)
var HTMLTemplate *template.Template

type InternalStorage struct {
	GaugeMetrics map[string]metrics.GaugeMetric
	CounterMetrics map[string]metrics.CounterMetric
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Config struct {
	Address string 				`env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile string 			`env:"STORE_FILE"`
	Restore bool 				`env:"RESTORE"`
}

func SetupRouters(r *chi.Mux) {
	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", GaugeMetricHandler)
		r.Post("/counter/{metricName}/{metricValue}", CounterMetricHandler)
		r.Post("/{metricName}/", NotFoundHandler)
		r.Post("/*", UnknownTypeMetricHandler)
		r.Post("/", UpdateMetricsJSONHandler)
	})
	r.Post("/value/", GetMetricsJSONHandler)
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
		Base:  metrics.Base{Name: metricName, Typename: metrics.Gauge},
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
		Base:  metrics.Base{Name: metricName, Typename: metrics.Counter},
		Value: counterMetricsStorage[metricName].Value + value}
	w.WriteHeader(http.StatusOK)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func UnknownTypeMetricHandler(w http.ResponseWriter, r *http.Request) {
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
		"gaugeMetrics":   gaugeMetricsStorage,
		"counterMetrics": counterMetricsStorage,
	}
	HTMLTemplate.Execute(w, renderData)
}

// UpdateMetricsJSONHandler Передача метрик на сервер
func UpdateMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := Metrics{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Print("type: ", m.MType, ", id: ", m.ID)

	updateMetricsStorage(m)

	w.WriteHeader(http.StatusOK)
}

func updateMetricsStorage(m Metrics) {
	switch m.MType {
	case metrics.Gauge:
		gaugeMetricsStorage[m.ID] = metrics.GaugeMetric{
			Base:  metrics.Base{Name: m.ID, Typename: metrics.Gauge},
			Value: *m.Value}
	case metrics.Counter:
		counterMetricsStorage[m.ID] = metrics.CounterMetric{
			Base:  metrics.Base{Name: m.ID, Typename: metrics.Counter},
			Value: counterMetricsStorage[m.ID].Value + *m.Delta}
	}
}

// GetMetricsJSONHandler Получение метрик с сервера
func GetMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := Metrics{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Print("type: ", m.MType, ", id: ", m.ID)

	switch m.MType {
	case metrics.Gauge:
		v := gaugeMetricsStorage[m.ID].Value
		m.Value = &v
	case metrics.Counter:
		v := counterMetricsStorage[m.ID].Value
		m.Delta = &v
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)

	w.WriteHeader(http.StatusOK)
}

func SaveMetrics(cfg Config) {
	//log.Println("Saving metrics to file")
	ticker := time.NewTicker(cfg.StoreInterval)
	for {
		<-ticker.C
		log.Println("Saving metrics to file")
		saveMetricsImpl(cfg)
	}
}

func saveMetricsImpl(cfg Config) {
	flags := os.O_WRONLY|os.O_CREATE

	f, err := os.OpenFile(cfg.StoreFile, flags, 0777) //0644
	if err != nil {
		log.Fatal("cannot open file for writing: ", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	internalStorage := InternalStorage{
		GaugeMetrics: gaugeMetricsStorage,
		CounterMetrics: counterMetricsStorage,
	}

	//if err := encoder.Encode(gaugeMetricsStorage); err != nil {
	//	log.Fatal("cannot encode gaugeMetricsStorage: ", err)
	//}
	//
	//if err = encoder.Encode(counterMetricsStorage); err != nil {
	//	log.Fatal("cannot encode counterMetricsStorage: ", err)
	//}

	if err := encoder.Encode(internalStorage); err != nil {
		log.Fatal("cannot encode internal metrics: ", err)
	}

	// dummy test save
	//f.Write([]byte(`{"id":"llvm","type":"gauge","value":10}`))
}