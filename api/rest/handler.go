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
			service.SaveMetric(ctx, counter)

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
			// Не обработал err
			service.SaveMetric(ctx, gauge)

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
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		//log.Printf("SpecificMetricHandler called. Requested metric type is %s, name is %s", metricType, metricName)
		log.Printf("Requested to return metric %s of type %s from storage", metricName, metricType)

		// TODO! инвертировать на не найдено + выход
		if m, found := service.GetMetric(metricName); found {
			w.WriteHeader(http.StatusOK)
			// TODO! switch по Type; добавить default Type
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

func UpdateMetricsJSONHandler(service server.Processor, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO! defer body.Close() в клиенте (агенте)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Request to update server's storage. Request body: %s", string(body))
		// TODO! использовать chi.middleware (для такого-то endpoint такой-то mw) R.Use(...); добавить проверку
		// TODO! выставить ближе к концу
		w.Header().Set("Content-Type", "application/json")

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

		service.SaveMetric(r.Context(), modelMetric)
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

func GetMetricsJSONHandler(service server.Processor, key string) http.HandlerFunc {
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

			//requestedMetric.Hash = requestedMetric.GenerateHash(key)

			if modelMetric, ok := service.GetMetric(requestedMetric.ID); ok {
				log.Println("Found metric in storage")
				responseMetric, err := models.FromModelMetrics(modelMetric)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}

				responseMetric.Hash = responseMetric.GenerateHash(key)

				data, err := json.Marshal(responseMetric)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string(data))
			} else {
				log.Println("Metric not found in storage")
				http.Error(w, "Metric not found in storage", http.StatusNotFound)
			}
		}
	}
}

func PingDBHandler(service server.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
		log.Println("`/ping` handler called")

		if err := service.Ping(); err != nil {
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

		log.Printf("/updates/ handler called. Request body: %s", string(body))

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

			service.SaveMetric(r.Context(), modelMetric)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}