package rest

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/go-chi/chi"
	"net/http"
)

func metricTypeValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r * http.Request) {
		metricType := chi.URLParam(r, "metricType")
		if err := model.ValidateType(metricType); err != nil {
			reason := fmt.Sprintf("Invalid url request: unknown metric type %s", metricType)
			http.Error(w, reason, http.StatusNotImplemented)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func metricNameValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r * http.Request) {
		if metricName := chi.URLParam(r, "metricName"); metricName == "" {
			reason := fmt.Sprint("Invalid url request: empty metric name. Maybe missing /{metric_name}/{metric_value} in the end of url?")
			http.Error(w, reason, http.StatusNotFound)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func metricValueValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r * http.Request) {
		if metricValue := chi.URLParam(r, "metricValue"); metricValue == "" {
			reason := fmt.Sprint("Invalid url request: empty metric value. Maybe missing /{metric_value} in the end of url?")
			http.Error(w, reason, http.StatusNotFound)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}