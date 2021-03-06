package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/go-chi/chi"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func URLMetricHandler(service server.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		switch metricType {
		case model.KCounter: {
			delta, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				reason := fmt.Sprintf("can't parse %s to int: %v", metricValue, err)
				http.Error(w, reason, http.StatusBadRequest)
				return
			}

			counter := model.NewCounterMetric(metricName, delta)

			if err := service.SaveMetric(ctx, counter); err != nil {
				reason := fmt.Sprintf("failed to store counter metric with name %s and value %s : %v", metricName, metricValue, err)
				http.Error(w, reason, http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		case model.KGauge: {
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				reason := fmt.Sprintf("can't parse %s to float: %v", metricValue, err)
				http.Error(w, reason, http.StatusBadRequest)
				return
			}

			gauge := model.NewGaugeMetric(metricName, value)

			if err := service.SaveMetric(ctx, gauge); err != nil {
				reason := fmt.Sprintf("failed to store gauge metric with name %s and value %s : %v", metricName, metricValue, err)
				http.Error(w, reason, http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		default:
			http.Error(w, "Unknown metric type", http.StatusNotImplemented)
			return
		}
	}
}

func AllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Dummy Page</title>
	</head>
	<body>
		Mock dummy page
	</body>
	</html>`

	t, err := template.New("getMetricList").Parse(htmlTemplate)
	if err != nil {
		errCode := http.StatusInternalServerError
		http.Error(w, err.Error(), errCode)
		return
	}

	data := "AC/DC"

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_ = t.Execute(w, data)
}

func SpecificMetricHandler(service server.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		m, err := service.GetMetric(ctx, metricName)
		if err != nil {
			reason := fmt.Sprintf("Metric %s with type %s not found in storage : %v", metricName, metricType, err)
			http.Error(w, reason, http.StatusNotFound)
			return
		}

		switch m.Type {
		case model.KCounter: {
			delta := m.MustGetInt()
			if _, err := w.Write([]byte(fmt.Sprint(delta))); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Printf("Returned counter metric from storage with delta %v", delta)
			w.WriteHeader(http.StatusOK)
			return
		}

		case model.KGauge: {
			value := m.MustGetFloat()
			if _, err := w.Write([]byte(fmt.Sprint(value))); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Printf("Returned gauge metric from storage with value %v", value)
			w.WriteHeader(http.StatusOK)
			return
		}

		default:
			reason := fmt.Sprintf("Unknown metric named %s with type %s", m.Name, m.Type)
			http.Error(w, reason, http.StatusInternalServerError)
			return
		}
	}
}

func UpdateMetricsJSONHandler(service server.Processor, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var m models.Metrics
		if err := json.Unmarshal(body, &m); err != nil {
			log.Printf("Failed to unmarshal following request body: %s", string(body))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := m.Validate(); err != nil {
			log.Printf("Invalid metric in incoming request. Type %s is not implemented", m.MType)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		if err := m.ValidateHash(key); err != nil {
			log.Println("Hash mismatched")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		modelMetric, err := m.ToModelMetric()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := service.SaveMetric(r.Context(), modelMetric); err != nil {
			log.Printf("Failed to save metric to storage: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(m); err != nil {
			log.Println("Failed to encode metric from storage")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("Metric updated")
		w.Header().Set("Content-Type", "application/json")
	}
}

func GetMetricsJSONHandler(service server.Processor, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var requestedMetric models.Metrics
		if err := json.Unmarshal(body, &requestedMetric); err != nil {
			log.Printf("Failed to unmarshal following request body: %s", string(body))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		modelMetric, err := service.GetMetric(ctx, requestedMetric.ID)
		if err != nil {
			reason := fmt.Sprintf("Metric not found in storage: %v", err)
			log.Println(reason)
			http.Error(w, reason, http.StatusNotFound)
			return
		}

		responseMetric, err := models.NewMetricFromCanonical(modelMetric)
		if err != nil {
			reason := fmt.Sprintf("Failed to convert from model type into api type: %v", err)
			log.Println(reason)
			http.Error(w, reason, http.StatusNotFound)
			return
		}

		responseMetric.Hash = responseMetric.GenerateHash(key)
		if requestedMetric.Hash != "" {
			if requestedMetric.Hash != responseMetric.Hash {
				reason := fmt.Sprintf("Hash mismatched. Expected %s, got %s", responseMetric.Hash, requestedMetric.Hash)
				log.Println(reason)
				http.Error(w, reason, http.StatusBadRequest)
				return
			}
		}

		data, err := json.Marshal(responseMetric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprint(w, string(data)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func PingDBHandler(service server.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
		ctx := r.Context()

		if err := service.Ping(ctx); err != nil {
			http.Error(w, "Failed to ping database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func UpdatesMetricsJSONHandler(service server.Processor, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var metrics []models.Metrics
		err = json.Unmarshal(body, &metrics)
		if err != nil {
			log.Printf("Failed to unmarshal following request body: %s", string(body))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, m := range metrics {
			if err := m.Validate(); err != nil {
				log.Printf("Invalid metric in incoming request. Type %s is not implemented", m.MType)
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}

			if err := m.ValidateHash(key); err != nil {
				log.Println("Hash mismatched")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			modelMetric, err := m.ToModelMetric()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := service.SaveMetric(r.Context(), modelMetric); err != nil {
				log.Printf("Failed to save metric to storage: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(w).Encode(m); err != nil {
				log.Println("Failed to encode metric from storage")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}