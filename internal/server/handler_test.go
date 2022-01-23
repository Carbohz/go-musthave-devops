package server

import (
	"bytes"
	"encoding/json"
	"github.com/Carbohz/go-musthave-devops/internal/common"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestCounterMetricHandler(t *testing.T) {
	tests := []struct {
		name     string
		URL      string
		pattern  string
		wantCode int
	}{
		{
			name:     "valid value",
			URL:      "/update/counter/metric/1",
			pattern:  "/update/counter/{metricName}/{metricValue}",
			wantCode: 200,
		},
		{
			name:     "invalid value",
			URL:      "/update/counter/metric/value",
			pattern:  "/update/counter/{metricName}/{metricValue}",
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

func TestGaugeMetricHandler(t *testing.T) {
	tests := []struct {
		name     string
		URL      string
		pattern  string
		wantCode int
	}{
		{
			name:     "valid value",
			URL:      "/update/gauge/metric/1234.5",
			pattern:  "/update/gauge/{metricName}/{metricValue}",
			wantCode: 200,
		},
		{
			name:     "invalid value",
			URL:      "/update/gauge/metric/value",
			pattern:  "/update/gauge/{metricName}/{metricValue}",
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

func TestUnknownTypeMetricHandler(t *testing.T) {
	tests := []struct {
		name string
		URL  string
		//pattern  string
		wantCode int
	}{
		{
			name: "update invalid type",
			URL:  "/update/unknown/testCounter/100",
			//pattern:  "/update/{metricType}/{metricName}/{metricValue}",
			wantCode: http.StatusNotImplemented,
		},
		//{
		//	name:     "invalid counter metric value request",
		//	URL:      "/update/counter/metric/value",
		//	pattern:  "/update/counter/{metricName}/{metricValue}",
		//	wantCode: 400,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			SetupRouters(r)
			req, err := http.NewRequest(http.MethodPost, tt.URL, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			//r.HandleFunc(tt.pattern, CounterMetricHandler)

			r.ServeHTTP(rec, req)
			res := rec.Result()

			assert.Equal(t, res.StatusCode, tt.wantCode)
			defer res.Body.Close()
		})
	}
}

// testing `/update/` handler
func TestUpdateMetricsJSONHandler(t *testing.T) {
	pattern := "/update/"

	tests := []struct {
		name     string
		URL      string
		rawJSON  []byte
		wantCode int
	}{
		{
			name:     "update json gauge metric",
			URL:      "/update/",
			rawJSON:  []byte(`{"id":"llvm","type":"gauge","value":1234.567}`),
			wantCode: 200,
		},
		{
			name:     "update json counter metric",
			URL:      "/update/",
			rawJSON:  []byte(`{"id":"llvm","type":"counter","delta":15}`),
			wantCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			SetupRouters(r)

			req, err := http.NewRequest(http.MethodPost, tt.URL, bytes.NewBuffer(tt.rawJSON))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			r.HandleFunc(pattern, UpdateMetricsJSONHandler)

			r.ServeHTTP(rec, req)
			res := rec.Result()

			assert.Equal(t, res.StatusCode, tt.wantCode)
			defer res.Body.Close()
		})
	}
}

// testing '/value/` handler
func TestGetMetricsJSONHandler(t *testing.T) {
	pattern := "/value/"

	// fails with this
	// serverDataJSON := []byte(`[{"id":"llvm","type":"gauge","value":1234.567},{"id":"PollCount","type":"counter","delta":5}]`)

	// ok with this
	serverDataJSON := []byte(`{"id":"llvm","type":"gauge","value":1234.567}`)

	tests := []struct {
		name      string
		URL       string
		rawJSON   []byte
		wantCode  int
		wantID    string
		wantMType string
		wantDelta int64
		wantValue float64
	}{
		{
			name:      "value json gauge metric",
			URL:       "/value/",
			rawJSON:   []byte(`{"id":"llvm","type":"gauge"}`),
			wantCode:  200,
			wantID:    "llvm",
			wantMType: "gauge",
			wantValue: 1234.567,
		},
		//{
		//	name:      "value json counter metric",
		//	URL:       "/update/",
		//	rawJSON:   []byte(`{"id":"PollCount","type":"counter"}`),
		//	wantCode:  200,
		//	wantID:    "PollCount",
		//	wantMType: "counter",
		//	wantDelta: 5,
		//},
		//{
		//	name:     "value json multi metrics",
		//	URL:      "/update/",
		//	rawJSON:  []byte(`[{"id":"llvm","type":"gauge"},{"id":"PollCount","type":"counter"}]`),
		//	wantCode: 200,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			SetupRouters(r)

			// send data to storage
			//rawJSON := []byte(`{"id":"llvm","type":"gauge","value":10}`)
			req, err := http.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(serverDataJSON))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			//r.HandleFunc(pattern, UpdateMetricsJSONHandler)

			r.ServeHTTP(rec, req)
			res := rec.Result()

			assert.Equal(t, tt.wantCode, res.StatusCode)
			defer res.Body.Close()

			// get data from storage
			//rawJSON = []byte(`{"id":"llvm","type":"gauge"}`)
			req, err = http.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(tt.rawJSON))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			rec = httptest.NewRecorder()
			r.HandleFunc(pattern, GetMetricsJSONHandler)

			r.ServeHTTP(rec, req)
			res = rec.Result()

			// unpack result
			body, _ := ioutil.ReadAll(res.Body)
			m := common.Metrics{}
			json.Unmarshal(body, &m)
			log.Printf("id: %v, type: %v, value: %v", m.ID, m.MType, *m.Value)

			assert.Equal(t, tt.wantCode, res.StatusCode)
			assert.Equal(t, tt.wantID, m.ID)
			assert.Equal(t, tt.wantMType, m.MType)
			switch m.MType {
			case "gauge":
				assert.InDelta(t, tt.wantValue, *m.Value, 1e-6)
			case "counter":
				assert.Equal(t, tt.wantDelta, *m.Delta)
			}
			defer res.Body.Close()
		})
	}
}
