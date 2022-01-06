package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"github.com/Carbohz/go-musthave-devops/internal/server"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/Carbohz/go-musthave-devops/internal/common"
	"github.com/go-chi/chi"
)

var gaugeMetricsStorage = make(map[string]metrics.GaugeMetric)
var counterMetricsStorage = make(map[string]metrics.CounterMetric)
var HTMLTemplate *template.Template
var secretKey string

type InternalStorage struct {
	GaugeMetrics   map[string]metrics.GaugeMetric
	CounterMetrics map[string]metrics.CounterMetric
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
	w.Header().Set("Content-Type", "text/html")
	renderData := map[string]interface{}{
		"gaugeMetrics":   gaugeMetricsStorage,
		"counterMetrics": counterMetricsStorage,
	}
	HTMLTemplate.Execute(w, renderData)
}

//// UpdateMetricsJSONHandler Передача метрик на сервер /update/
//func UpdateMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//
//	log.Printf("/update/ handler called. Request body: %s", string(body))
//
//	m := common.Metrics{}
//	err = json.Unmarshal(body, &m)
//	if err != nil {
//		log.Printf("Failed to unmarshal following request body: %s", string(body))
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	err = m.CheckHash(secretKey)
//	if err == nil {
//		log.Println("Hash matched, updating internal server metrics")
//		updateMetricsStorage(m)
//		w.WriteHeader(http.StatusOK)
//		return
//	} else {
//		log.Println("Hash mismatched, bad request")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//}

// UpdateMetricsJSONHandler Передача метрик на сервер
func UpdateMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := common.Metrics{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	updateMetricsStorage(m)

	w.WriteHeader(http.StatusOK)
}

func updateMetricsStorage(m common.Metrics) {
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

// GetMetricsJSONHandler Получение метрик с сервера /value/
func GetMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("/value/ handler called. Request body: %s", string(body))

	w.Header().Set("Content-Type", "application/json")
	if string(body)[0] == '[' {
		log.Println("Request body contains array of metrics")
		json.NewEncoder(w).Encode(generateMultipleMetrics(body))
	} else {
		log.Println("Request body contains single metric")
		json.NewEncoder(w).Encode(generateSingleMetric(body))
	}
}

func DumpMetrics(cfg server.Config) {
	ticker := time.NewTicker(cfg.StoreInterval)
	for {
		<-ticker.C
		log.Printf("Dumping metrics to file %s", cfg.StoreFile)
		DumpMetricsImpl(cfg)
	}
}

func DumpMetricsImpl(cfg server.Config) {
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(cfg.StoreFile, flag, 0644)
	if err != nil {
		log.Fatal("Can't open file for dumping: ", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	internalStorage := InternalStorage{
		GaugeMetrics:   gaugeMetricsStorage,
		CounterMetrics: counterMetricsStorage,
	}

	if err := encoder.Encode(internalStorage); err != nil {
		log.Fatal("Can't encode server's metrics: ", err)
	}
}

func LoadMetrics(cfg server.Config) {
	log.Printf("Loading metrics from file %s", cfg.StoreFile)

	flag := os.O_RDONLY

	f, err := os.OpenFile(cfg.StoreFile, flag, 0)
	if err != nil {
		log.Print("Can't open file for loading metrics: ", err)
		return
	}
	defer f.Close()

	var internalStorage InternalStorage

	if err := json.NewDecoder(f).Decode(&internalStorage); err != nil {
		log.Fatal("Can't decode metrics: ", err)
	}

	gaugeMetricsStorage = internalStorage.GaugeMetrics
	counterMetricsStorage = internalStorage.CounterMetrics
	log.Printf("Metrics successfully loaded from file %s", cfg.StoreFile)
}

func PassSecretKey(key string) {
	secretKey = key
}

func generateSingleMetric(body []byte) common.Metrics {
	m := common.Metrics{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		log.Println("`GetMetricsJSONHandler` error triggered - `/value/` handler")
		log.Println("Single metric json body error!")
		log.Printf("Unmarshalling JSON error: %v", err)
		log.Printf("Request body was: %s", string(body))
		//http.Error(w, err.Error(), http.StatusBadRequest)
	}

	log.Println("***Initially value-requested single metric***")
	switch m.MType {
	case "gauge":
		log.Printf("[ID: %v, Type: %v]", m.ID, m.MType)
	case "counter":
		log.Printf("[ID: %v, Type: %v]", m.ID, m.MType)
	default:
		log.Println("Unknown metric type")
	}

	switch m.MType {
	case metrics.Gauge:
		v := gaugeMetricsStorage[m.ID].Value
		m.Value = &v
	case metrics.Counter:
		v := counterMetricsStorage[m.ID].Value
		m.Delta = &v
	}

	log.Println("***Filled with server values single metric***")
	switch m.MType {
	case "gauge":
		log.Printf("[ID: %v, Type: %v, Value: %v]", m.ID, m.MType, *m.Value)
	case "counter":
		log.Printf("[ID: %v, Type: %v, Value: %v]", m.ID, m.MType, *m.Delta)
	default:
		log.Println("Unknown metric type")
	}
	return m
}

func generateMultipleMetrics(body []byte) []common.Metrics {
	var mArr []common.Metrics
	err := json.Unmarshal(body, &mArr)
	if err != nil {
		log.Println("`GetMetricsJSONHandler` error triggered - `/value/` handler")
		log.Println("Array of metrics json body error!")
		log.Printf("Unmarshalling JSON error: %v", err)
		log.Printf("Request body was: %s", string(body))
	}

	log.Println("***Initially value-requested metrics array***")
	for _, mtrc := range mArr {
		switch mtrc.MType {
		case "gauge":
			log.Printf("[ID: %v, Type: %v]", mtrc.ID, mtrc.MType)
		case "counter":
			log.Printf("[ID: %v, Type: %v]", mtrc.ID, mtrc.MType)
		default:
			log.Println("Unknown metric type")
		}
	}

	for i, m := range mArr {
		switch m.MType {
		case metrics.Gauge:
			v := gaugeMetricsStorage[m.ID].Value
			// m.Value = &v
			mArr[i].Value = &v
		case metrics.Counter:
			v := counterMetricsStorage[m.ID].Value
			// m.Delta = &v
			mArr[i].Delta = &v
		}
	}

	log.Println("***Filled with server values metrics array***")
	for _, mtrc := range mArr {
		switch mtrc.MType {
		case "gauge":
			log.Printf("[ID: %v, Type: %v, Value: %v]", mtrc.ID, mtrc.MType, *mtrc.Value)
		case "counter":
			log.Printf("[ID: %v, Type: %v, Value: %v]", mtrc.ID, mtrc.MType, *mtrc.Delta)
		default:
			log.Println("Unknown metric type")
		}
	}
	return mArr
}
