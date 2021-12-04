package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeMetricHandler(t *testing.T) {
	tests := []struct {
		name       string
		requestURL string
		wantCode   int
	}{
		{
			name:       "valid gauge metric value request",
			requestURL: "/update/gauge/metric/1234.5",
			wantCode:   200,
		},
		{
			name:       "invalid gauge metric value request",
			requestURL: "/update/gauge/metric/value",
			wantCode:   400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.requestURL, nil)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(GaugeMetricHandler)

			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.wantCode)
			defer res.Body.Close()
		})
	}
}

func TestCounterMetricHandler(t *testing.T) {
	tests := []struct {
		name       string
		requestURL string
		wantCode   int
	}{
		{
			name:       "valid counter metric value request",
			requestURL: "/update/counter/metric/1",
			wantCode:   200,
		},
		{
			name:       "invalid counter metric value request",
			requestURL: "/update/counter/metric/value",
			wantCode:   400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.requestURL, nil)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(CounterMetricHandler)

			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.wantCode)
			defer res.Body.Close()
		})
	}
}
