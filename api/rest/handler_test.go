package rest

import (
	"bytes"
	"encoding/json"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/Carbohz/go-musthave-devops/model"
	v1 "github.com/Carbohz/go-musthave-devops/service/server/v1"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/markphelps/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	storagemock "github.com/Carbohz/go-musthave-devops/storage/mock"
)

func TestUpdateMetricWithURL(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "Valid gauge metric1",
			path: "/update/gauge/metric1/123.45",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Valid counter metric2",
			path: "/update/counter/metric2/123",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MType",
			path: "/update/abcdef/metric3/123",
			want: want{
				code: http.StatusNotImplemented,
			},
		},
	}

	mockCtrl := gomock.NewController(t)

	metricStorage := storagemock.NewMockMetricsStorager(mockCtrl)
	processor, _ := v1.NewService(metricStorage)
	r := chi.NewRouter()
	setupRouters(r, processor)

	server := httptest.NewServer(r)
	defer server.Close()

	metric1 := model.Metric{Name: "metric1", Type: model.KGauge, Value: optional.NewFloat64(123.45)}
	metric2 := model.Metric{Name: "metric2", Type: model.KCounter, Delta: optional.NewInt64(123)}

	gomock.InOrder(
		metricStorage.EXPECT().SaveMetric(metric1),
		metricStorage.EXPECT().SaveMetric(metric2),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _ := helperDoRequest(t, server, http.MethodPost, tt.path, nil)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestUpdateMetricWithBody(t *testing.T) {
	metric1Value := 123.45
	var metric2Delta int64 = 123

	type want struct {
		code int
	}
	tests := []struct {
		name   string
		path   string
		metric models.Metrics
		want   want
	}{
		{
			name:   "Valid gauge metric1",
			path:   "/update/",
			//metric: model.Metric{Name: "metric1", Type: model.KGauge, Value: optional.NewFloat64(123.45)},
			metric: models.Metrics{ID: "metric1", MType: model.KGauge, Value: &metric1Value},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "Valid counter metric2",
			path:   "/update/",
			//metric: model.Metric{Name: "metric2", Type: model.KCounter, Delta: optional.NewInt64(123)},
			metric: models.Metrics{ID: "metric2", MType: model.KCounter, Delta: &metric2Delta},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MType",
			path: "/update/",
			//metric: model.Metric{ID: "metric3", MType: model.MetricType("abcdef")},
			//metric: model.Metric{Name: "metric3", Type: "abcdef"},
			metric: models.Metrics{ID: "metric3", MType: "abcdef"},
			want: want{
				code: http.StatusNotImplemented,
			},
		},
	}

	mockCtrl := gomock.NewController(t)

	metricStorage := storagemock.NewMockMetricsStorager(mockCtrl)
	metric1 := model.Metric{Name: "metric1", Type: model.KGauge, Value: optional.NewFloat64(123.45)}
	metric2 := model.Metric{Name: "metric2", Type: model.KCounter, Delta: optional.NewInt64(123)}
	//metric3 := model.Metric{Name: "metric3", Type: "abcdef"}

	gomock.InOrder(
		metricStorage.EXPECT().SaveMetric(metric1),
		metricStorage.EXPECT().SaveMetric(metric2),
		//metricStorage.EXPECT().SaveMetric(metric3),
	)
	processor, _ := v1.NewService(metricStorage)
	r := chi.NewRouter()
	setupRouters(r, processor)

	server := httptest.NewServer(r)
	defer server.Close()


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.metric)
			log.Printf("Marshalled data: %s", string(data))
			require.NoError(t, err)
			statusCode, _ := helperDoRequest(t, server, http.MethodPost, tt.path, &data)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestGetMetricWithBody(t *testing.T) {
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		path   string
		metric models.Metrics
		want   want
	}{
		{
			name:   "Valid gauge metric1",
			path:   "/value/",
			metric: models.Metrics{ID: "metric1", MType: model.KGauge},
			want: want{
				code: http.StatusOK,
				body: "{\"id\":\"metric1\",\"type\":\"gauge\",\"value\":123.45}",
			},
		},
		{
			name:   "Valid counter metric2",
			path:   "/value/",
			metric: models.Metrics{ID: "metric2", MType: model.KCounter},
			want: want{
				code: http.StatusOK,
				body: "{\"id\":\"metric2\",\"type\":\"counter\",\"delta\":123}",
			},
		},
		{
			name:   "Invalid MType",
			path:   "/value/",
			metric: models.Metrics{ID: "metric3", MType: "abrakadabra"},
			want: want{
				//code: http.StatusNotImplemented,
				code: http.StatusNotFound,
				body: "Metric not found in storage\n",
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	metricStorage := storagemock.NewMockMetricsStorager(mockCtrl)
	processor, _ := v1.NewService(metricStorage)
	r := chi.NewRouter()
	setupRouters(r, processor)

	server := httptest.NewServer(r)
	defer server.Close()

	metric1 := model.Metric{Name: "metric1", Type: model.KGauge, Value: optional.NewFloat64(123.45)}
	metric2 := model.Metric{Name: "metric2", Type: model.KCounter, Delta: optional.NewInt64(123)}
	metric3 := model.Metric{Name: "metric3", Type: "abrakadabra"}

	gomock.InOrder(
		metricStorage.EXPECT().GetMetric(gomock.Any()).Return(metric1, true),
		metricStorage.EXPECT().GetMetric(gomock.Any()).Return(metric2, true),
		metricStorage.EXPECT().GetMetric(gomock.Any()).Return(metric3, false),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.metric)
			require.NoError(t, err)
			statusCode, body := helperDoRequest(t, server, http.MethodPost, tt.path, &data)
			assert.Equal(t, tt.want.code, statusCode)
			assert.Equal(t, tt.want.body, body)
		})
	}
}

func helperDoRequest(t *testing.T, server *httptest.Server, method, path string, data *[]byte) (int, string) {
	var body io.Reader
	if data != nil {
		body = bytes.NewBuffer(*data)
	}
	request, err := http.NewRequest(method, server.URL+path, body)
	require.NoError(t, err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	responseBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	defer response.Body.Close()

	return response.StatusCode, string(responseBody)
}
