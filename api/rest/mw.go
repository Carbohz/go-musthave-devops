package rest

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/go-chi/chi"
	"net/http"
)

func urlValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		if err := model.ValidateType(metricType); err != nil {
			reason := fmt.Sprintf("Invalid url request: unknown metric type %s", metricType)
			http.Error(w, reason, http.StatusNotImplemented)
			return
		}

		if metricName := chi.URLParam(r, "metricName"); metricName == "" {
			reason := fmt.Sprintf("Invalid url request: missing /{metric_name}/{metric_value} in the end of url")
			http.Error(w, reason, http.StatusNotFound)
		}

		if metricValue := chi.URLParam(r, "metricValue"); metricValue == "" {
			reason := fmt.Sprintf("Invalid url request: missing /{metric_value} in the end of url")
			http.Error(w, reason, http.StatusNotFound)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}