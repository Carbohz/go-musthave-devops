package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/agent"
	"github.com/Carbohz/go-musthave-devops/internal/common"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"io"
	"log"
	"net/http"
)

type AllMetrics struct {
	RuntimeMetrics    []metrics.GaugeMetric
	RandomValueMetric metrics.GaugeMetric
	PollCountMetric   metrics.CounterMetric
	Hash              string `json:"hash,omitempty"` // значение хеш-функции
}

func SendGaugeMetric(client *http.Client, m metrics.GaugeMetric, address string) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%f", address, m.Typename, m.Name, m.Value)
	return SendMetric(client, url, m.Base)
}

func SendCounterMetric(client *http.Client, m metrics.CounterMetric, address string) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%d", address, m.Typename, m.Name, m.Value)
	return SendMetric(client, url, m.Base)
}

func SendMetric(client *http.Client, url string, m metrics.Base) error {
	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		log.Printf("Failed to \"Post\" request to update metric \"%s\" of type \"%s\"", m.Name, m.Typename)
		log.Printf("Error: %v", err)
		return err
	}
	defer resp.Body.Close()
	return err
}

func SendMetricsJSON(client *http.Client, runtimeMetrics []metrics.GaugeMetric,
	randomValueMetric metrics.GaugeMetric,
	pollCountMetric metrics.CounterMetric, cfg agent.Config) error {
	//url := fmt.Sprintf("http://%s/value/", cfg.Address)
	url := fmt.Sprintf("http://%s/update/", cfg.Address)
	log.Printf("Sending JSON metrics to url: %s", url)

	allMetrics := createMetricsArr(runtimeMetrics, randomValueMetric, pollCountMetric, cfg.Key)

	body := bytes.NewBuffer(marshallMetricsJSON(allMetrics))

	log.Printf("Request body (JSON): %v", body)

	resp, err := client.Post(url, "application/json", body)
	if err != nil {
		log.Printf("Failed to \"Post\" metrics in JSON format. Error: %v", err)
		return err
	}

	log.Printf("Response: %v", resp)

	defer resp.Body.Close()
	return err
}

func marshallMetricsJSON(allMetrics []common.Metrics) []byte {
	rawJSON, err := json.Marshal(allMetrics)
	if err != nil {
		log.Fatalf("Error occured during metrics marshalling: %v", err)
	}
	// log.Printf("Generated raw JSON: %v", string(rawJSON))
	return rawJSON
}

func marshallMetricJSON(m common.Metrics) []byte {
	rawJSON, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Error occured during metrics marshalling: %v", err)
	}
	// log.Printf("Generated raw JSON: %v", string(rawJSON))
	return rawJSON
}

func createMetricsArr(runtimeMetrics []metrics.GaugeMetric,
	randomValueMetric metrics.GaugeMetric,
	pollCountMetric metrics.CounterMetric, key string) []common.Metrics {
	var metricsArr []common.Metrics

	metricsArr = append(metricsArr, generateCommonRuntimeMetrics(runtimeMetrics, key)...)

	metricsArr = append(metricsArr, generateCommonRandomValueMetric(randomValueMetric, key))

	metricsArr = append(metricsArr, generateCommonPollCountMetric(pollCountMetric, key))

	return metricsArr
}

func generateCommonRuntimeMetrics(runtimeMetrics []metrics.GaugeMetric, key string) []common.Metrics {
	var metricsArr []common.Metrics

	for _, m := range runtimeMetrics {
		value := m.Value
		currentMetric := common.Metrics{ID: m.Name, MType: m.Typename, Delta: nil, Value: &value, Hash: ""}
		currentMetric.Hash = currentMetric.GenerateHash(key)

		metricsArr = append(metricsArr, currentMetric)
	}

	return metricsArr
}

func generateCommonRandomValueMetric(m metrics.GaugeMetric, key string) common.Metrics {
	value := m.Value
	currentMetric := common.Metrics{ID: m.Name, MType: m.Typename, Delta: nil, Value: &value, Hash: ""}
	currentMetric.Hash = currentMetric.GenerateHash(key)

	return currentMetric
}

func generateCommonPollCountMetric(m metrics.CounterMetric, key string) common.Metrics {
	delta := m.Value
	currentMetric := common.Metrics{ID: m.Name, MType: m.Typename, Delta: &delta, Value: nil, Hash: ""}
	currentMetric.Hash = currentMetric.GenerateHash(key)

	return currentMetric
}

func generateCommonGaugeMetric(m metrics.GaugeMetric, key string) common.Metrics {
	value := m.Value
	currentMetric := common.Metrics{ID: m.Name, MType: m.Typename, Delta: nil, Value: &value, Hash: ""}
	if key != "" {
		currentMetric.Hash = currentMetric.GenerateHash(key)
	}
	return currentMetric
}

func generateCommonCounterMetric(m metrics.CounterMetric, key string) common.Metrics {
	delta := m.Value
	currentMetric := common.Metrics{ID: m.Name, MType: m.Typename, Delta: &delta, Value: nil, Hash: ""}
	if key != "" {
		currentMetric.Hash = currentMetric.GenerateHash(key)
	}
	return currentMetric
}

func SendGaugeMetricJSON(client *http.Client, m metrics.GaugeMetric, cfg agent.Config) error {
	url := fmt.Sprintf("http://%s/update/", cfg.Address)
	log.Printf("Sending JSON metric to url: %s", url)

	commonGaugeMetric := generateCommonGaugeMetric(m, cfg.Key)

	bodyRaw := marshallMetricJSON(commonGaugeMetric)

	body := bytes.NewBuffer(bodyRaw)

	log.Printf("Request body (JSON): %v", body)

	return SendMetricJSON(client, url, body)
}

func SendCounterMetricJSON(client *http.Client, m metrics.CounterMetric, cfg agent.Config) error {
	url := fmt.Sprintf("http://%s/update/", cfg.Address)
	log.Printf("Sending JSON metric to url: %s", url)

	commonCounterMetric := generateCommonCounterMetric(m, cfg.Key)

	bodyRaw := marshallMetricJSON(commonCounterMetric)

	body := bytes.NewBuffer(bodyRaw)

	log.Printf("Request body (JSON): %v", body)

	return SendMetricJSON(client, url, body)
}

func SendMetricJSON(client *http.Client, url string, body io.Reader) error {
	resp, err := client.Post(url, "application/json", body)
	if err != nil {
		log.Printf("Failed to \"Post\" metrics in JSON format. Error: %v", err)
		return err
	}

	log.Printf("Response: %v", *resp)

	defer resp.Body.Close()
	return err
}
