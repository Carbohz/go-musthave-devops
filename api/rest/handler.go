package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/markphelps/optional"
	"io/ioutil"
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

func UpdateMetricsJSONHandler(service server.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Request to update server's storage. Request body: %s", string(body))
		w.Header().Set("Content-Type", "application/json")

		var m models.Metrics
		err = json.Unmarshal(body, &m)
		if err != nil {
			log.Printf("Failed to unmarshal following request body: %s", string(body))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		modelMetric := m.ToModelMetric()
		service.ProcessMetric(r.Context(), modelMetric)
		err = json.NewEncoder(w).Encode(m)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			log.Printf("Failed to update metric on storage")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			log.Println("Metric updated")
		}
	}
}

func GetMetricsJSONHandler(service server.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Request to return metric from storage. Request body: %s", string(body))

		w.Header().Set("Content-Type", "application/json")
		if string(body)[0] == '[' {
			log.Println("Request body contains array of metrics. Currently not supported")
			http.Error(w, "Request body contains array of metrics. Currently not supported", http.StatusBadRequest)
		} else {
			log.Println("Request body contains single metric")

			var requestedMetric models.Metrics
			if err := json.Unmarshal(body, &requestedMetric); err != nil {
				log.Printf("Failed to unmarshal following request body: %s", string(body))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if modelMetric, ok := service.GetMetric(requestedMetric.ID); ok {
				log.Println("Found metric in storage")
				responseMetric := models.FromModelMetrics(modelMetric)
				json.NewEncoder(w).Encode(responseMetric)
			} else {
				log.Println("Metric not found in storage")
			}
		}
	}
}