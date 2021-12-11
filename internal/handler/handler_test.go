package handler

import (
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeMetricHandler(t *testing.T) {
	tests := []struct {
		name       string
		URL string
		pattern string
		wantCode   int
	}{
		{
			name: "valid gauge metric value request",
			URL: "/update/gauge/metric/1234.5",
			pattern: "/update/gauge/{metricName}/{metricValue}",
			wantCode: 200,
		},
		{
			name: "invalid gauge metric value request",
			URL: "/update/gauge/metric/value",
			pattern: "/update/gauge/{metricName}/{metricValue}",
			wantCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			SetupRouters(r)
			req, err := http.NewRequest(http.MethodPost, tt.URL, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			r.HandleFunc(tt.pattern, GaugeMetricHandler)

			r.ServeHTTP(rec, req)
			res := rec.Result()

			assert.Equal(t, res.StatusCode, tt.wantCode)
			defer res.Body.Close()
		})
	}
}

func TestCounterMetricHandler(t *testing.T) {
	tests := []struct {
		name string
		URL string
		pattern string
		wantCode int
	}{
		{
			name: "valid counter metric value request",
			URL: "/update/counter/metric/1",
			pattern: "/update/counter/{metricName}/{metricValue}",
			wantCode: 200,
		},
		{
			name: "invalid counter metric value request",
			URL: "/update/counter/metric/value",
			pattern: "/update/counter/{metricName}/{metricValue}",
			wantCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			SetupRouters(r)
			req, err := http.NewRequest(http.MethodPost, tt.URL, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			r.HandleFunc(tt.pattern, CounterMetricHandler)

			r.ServeHTTP(rec, req)
			res := rec.Result()

			assert.Equal(t, res.StatusCode, tt.wantCode)
			defer res.Body.Close()
		})
	}
}
