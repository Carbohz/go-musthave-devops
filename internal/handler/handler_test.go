package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeHandler(t *testing.T) {
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

func TestCounterHandler(t *testing.T) {
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

func TestGetRequestBody(t *testing.T) {
	tests := []struct {
		name       string
		requestURL string
		wantToken1 string
		wantToken2 string
	}{
		{
			name: "gauge request tokens",
			requestURL: "/update/gauge/Alloc/111",
			wantToken1: "Alloc",
			wantToken2: "111",
		},
		{
			name: "counter request tokens",
			requestURL: "/update/counter/PollCount/222",
			wantToken1: "PollCount",
			wantToken2: "222",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.requestURL, nil)

			gotToken1, gotToken2 := GetRequestBody(request)
			assert.Equal(t, gotToken1, tt.wantToken1)
			assert.Equal(t, gotToken2, tt.wantToken2)
		})
	}
}
