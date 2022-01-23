package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"strconv"
	"text/template"
	"time"

	"github.com/Carbohz/go-musthave-devops/internal/common"
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//var gaugeMetricsStorage = make(map[string]metrics.GaugeMetric)
//var counterMetricsStorage = make(map[string]metrics.CounterMetric)
var HTMLTemplate *template.Template
var secretKey string
var db *sql.DB

var instance Instance

//type InternalStorage struct {
//	GaugeMetrics   map[string]metrics.GaugeMetric
//	CounterMetrics map[string]metrics.CounterMetric
//}

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
	r.Get("/ping", PingDBHandler)
}

func GaugeMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "parsing error. Bad request", http.StatusBadRequest)
		return
	}
	//gaugeMetricsStorage[metricName] = metrics.GaugeMetric{
	//	Base:  metrics.Base{Name: metricName, Typename: metrics.Gauge},
	//	Value: value}
	instance.StoreGaugeMetric(metricName, value)
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
	//counterMetricsStorage[metricName] = metrics.CounterMetric{
	//	Base:  metrics.Base{Name: metricName, Typename: metrics.Counter},
	//	Value: counterMetricsStorage[metricName].Value + value}
	instance.StoreCounterMetric(metricName, value)
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
		value, err := instance.FindCounterMetric(metricName)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(value)))
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return

		//if val, found := counterMetricsStorage[metricName]; found {
		//	w.WriteHeader(http.StatusOK)
		//	w.Write([]byte(fmt.Sprint(val.Value)))
		//	return
		//}
		//reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
		//http.Error(w, reason, http.StatusNotFound)
		//return
	}

	if metricType == metrics.Gauge {
		value, err := instance.FindCounterMetric(metricName)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(value)))
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return

		//if val, found := gaugeMetricsStorage[metricName]; found {
		//	w.WriteHeader(http.StatusOK)
		//	w.Write([]byte(fmt.Sprint(val.Value)))
		//	return
		//}
		//reason := fmt.Sprintf("Unknown metric \"%s\" of type \"%s\"", metricName, metricType)
		//http.Error(w, reason, http.StatusNotFound)
	}
}

func AllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//renderData := map[string]interface{}{
	//	"gaugeMetrics":   gaugeMetricsStorage,
	//	"counterMetrics": counterMetricsStorage,
	//}
	renderData := map[string]interface{}{
		"gaugeMetrics":   instance.GetGaugeMetrics(),
		"counterMetrics": instance.GetCounterMetrics(),
	}
	HTMLTemplate.Execute(w, renderData)
}

// UpdateMetricsJSONHandler Передача метрик на сервер /update/
func UpdateMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("/update/ handler called. Request body: %s", string(body))

	m := common.Metrics{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Printf("Failed to unmarshal following request body: %s", string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.CheckHash(secretKey)
	if err == nil {
		log.Println("Hash matched, updating internal server metrics")
		updateMetricsStorage(m)
		//w.WriteHeader(http.StatusOK)

		//response := m.Hash
		//log.Printf("Response message: %s", response)
		//w.Write([]byte(response))
		//w.Write(generateResponseJSON(m))
		err = json.NewEncoder(w).Encode(m)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			log.Printf("Error occurred during response json encoding: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	} else {
		log.Println("Hash mismatched, bad request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func updateMetricsStorage(m common.Metrics) {
	//switch m.MType {
	//case metrics.Gauge:
	//	gaugeMetricsStorage[m.ID] = metrics.GaugeMetric{
	//		Base:  metrics.Base{Name: m.ID, Typename: metrics.Gauge},
	//		Value: *m.Value}
	//case metrics.Counter:
	//	counterMetricsStorage[m.ID] = metrics.CounterMetric{
	//		Base:  metrics.Base{Name: m.ID, Typename: metrics.Counter},
	//		Value: counterMetricsStorage[m.ID].Value + *m.Delta}
	//}

	switch m.MType {
	case metrics.Gauge:
		instance.UpdateGaugeMetrics(m.ID, *m.Value)
	case metrics.Counter:
		instance.UpdateCounterMetrics(m.ID, *m.Delta)
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

//func DumpMetrics(cfg server.Config) {
//	ticker := time.NewTicker(cfg.StoreInterval)
//	for {
//		<-ticker.C
//		log.Printf("Dumping metrics to file %s", cfg.StoreFile)
//		DumpMetricsImpl(cfg)
//	}
//}

//func DumpMetricsImpl(cfg server.Config) {
//	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
//
//	f, err := os.OpenFile(cfg.StoreFile, flag, 0644)
//	if err != nil {
//		log.Fatal("Can't open file for dumping: ", err)
//	}
//	defer f.Close()
//
//	encoder := json.NewEncoder(f)
//
//	internalStorage := InternalStorage{
//		GaugeMetrics:   gaugeMetricsStorage,
//		CounterMetrics: counterMetricsStorage,
//	}
//
//	if err := encoder.Encode(internalStorage); err != nil {
//		log.Fatal("Can't encode server's metrics: ", err)
//	}
//}

//func LoadMetrics(cfg server.Config) {
//	log.Printf("Loading metrics from file %s", cfg.StoreFile)
//
//	flag := os.O_RDONLY
//
//	f, err := os.OpenFile(cfg.StoreFile, flag, 0)
//	if err != nil {
//		log.Print("Can't open file for loading metrics: ", err)
//		return
//	}
//	defer f.Close()
//
//	var internalStorage InternalStorage
//
//	if err := json.NewDecoder(f).Decode(&internalStorage); err != nil {
//		log.Fatal("Can't decode metrics: ", err)
//	}
//
//	gaugeMetricsStorage = internalStorage.GaugeMetrics
//	counterMetricsStorage = internalStorage.CounterMetrics
//	log.Printf("Metrics successfully loaded from file %s", cfg.StoreFile)
//}

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
	case metrics.Gauge:
		log.Printf("[ID: %v, Type: %v]", m.ID, m.MType)
	case metrics.Counter:
		log.Printf("[ID: %v, Type: %v]", m.ID, m.MType)
	default:
		log.Println("Unknown metric type")
	}

	switch m.MType {
	case metrics.Gauge:
		//v := gaugeMetricsStorage[m.ID].Value
		v, _ := instance.FindGaugeMetric(m.ID)
		m.Value = &v
	case metrics.Counter:
		//v := counterMetricsStorage[m.ID].Value
		v, _ := instance.FindCounterMetric(m.ID)
		m.Delta = &v
	}

	m.Hash = m.GenerateHash(secretKey)

	log.Println("***Filled with server values single metric***")
	switch m.MType {
	case metrics.Gauge:
		log.Printf("[ID: %v, Type: %v, Value: %v, Hash: %s]", m.ID, m.MType, *m.Value, m.Hash)
	case metrics.Counter:
		log.Printf("[ID: %v, Type: %v, Value: %v, Hash: %s]", m.ID, m.MType, *m.Delta, m.Hash)
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
		case metrics.Gauge:
			log.Printf("[ID: %v, Type: %v]", mtrc.ID, mtrc.MType)
		case metrics.Counter:
			log.Printf("[ID: %v, Type: %v]", mtrc.ID, mtrc.MType)
		default:
			log.Println("Unknown metric type")
		}
	}

	for i, m := range mArr {
		switch m.MType {
		case metrics.Gauge:
			//v := gaugeMetricsStorage[m.ID].Value
			v, _ := instance.FindGaugeMetric(m.ID)
			mArr[i].Value = &v
		case metrics.Counter:
			//v := counterMetricsStorage[m.ID].Value
			v, _ := instance.FindCounterMetric(m.ID)
			mArr[i].Delta = &v
		}

		mArr[i].GenerateHash(secretKey)
	}

	log.Println("***Filled with server values metrics array***")
	for _, mtrc := range mArr {
		switch mtrc.MType {
		case metrics.Gauge:
			log.Printf("[ID: %v, Type: %v, Value: %v, Hash: %s]", mtrc.ID, mtrc.MType, *mtrc.Value, mtrc.Hash)
		case metrics.Counter:
			log.Printf("[ID: %v, Type: %v, Value: %v, Hash: %s]", mtrc.ID, mtrc.MType, *mtrc.Delta, mtrc.Hash)
		default:
			log.Println("Unknown metric type")
		}
	}
	return mArr
}

func generateResponseJSON(m common.Metrics) []byte {
	rawJSON, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Error occured during metrics marshalling: %v", err)
	}
	log.Printf("Generated JSON response: %s", string(rawJSON))
	return rawJSON
}

func PingDBHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("`/ping` handler called")
	if db == nil {
		log.Printf("Connection error: database is not connected")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("200 OK"))
}

func ConnectDB(dbPath string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("pgx", dbPath)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}
	return db, nil
}
