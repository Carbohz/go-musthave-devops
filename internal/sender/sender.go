package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/agent"
	"github.com/Carbohz/go-musthave-devops/internal/common"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"log"
	"net/http"
)

type AllMetrics struct {
	RuntimeMetrics    []metrics.GaugeMetric
	RandomValueMetric metrics.GaugeMetric
	PollCountMetric   metrics.CounterMetric
	Hash              string `json:"hash,omitempty"` // значение хеш-функции
}

//var allMetricsss []common.Metrics

func SendGaugeMetric(client *http.Client, m metrics.GaugeMetric, address string) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%f", address, m.Typename, m.Name, m.Value)
	return Send(client, url, m.Base)
}

func SendCounterMetric(client *http.Client, m metrics.CounterMetric, address string) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%d", address, m.Typename, m.Name, m.Value)
	return Send(client, url, m.Base)
}

func Send(client *http.Client, url string, m metrics.Base) error {
	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		log.Printf("Failed to \"Post\" request to update metric \"%s\" of type \"%s\"", m.Name, m.Typename)
		return err
	}
	defer resp.Body.Close()
	return err
}

func SendMetricsJSON(client *http.Client, runtimeMetrics []metrics.GaugeMetric,
	randomValueMetric metrics.GaugeMetric,
	pollCountMetric metrics.CounterMetric, cfg agent.Config) error {
	url := fmt.Sprintf("http://%s/value/", cfg.Address)

	allMetrics := createMetricsArr(runtimeMetrics, randomValueMetric, pollCountMetric, cfg.Key)

	body := marshallMetricsJSON(allMetrics)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Failed to \"Post\" metrics in JSON format")
	}
	defer resp.Body.Close()
	return err
}

func marshallMetricsJSON(allMetrics []common.Metrics) []byte {
	rawJSON, err := json.Marshal(allMetrics)
	if err != nil {
		log.Fatalf("Error occured during metrics marshalling: %v", err)
	}
	fmt.Println(string(rawJSON))
	return rawJSON
}

func createMetricsArr(runtimeMetrics []metrics.GaugeMetric,
	randomValueMetric metrics.GaugeMetric,
	pollCountMetric metrics.CounterMetric, key string) []common.Metrics {
	var metricsArr []common.Metrics
	var currentMetric common.Metrics

	// add runtime metrics to array
	for _, m := range runtimeMetrics {
		currentMetric.ID = m.Name
		currentMetric.MType = m.Typename
		currentMetric.Value = &m.Value
		//currentMetric.Delta
		currentMetric.Hash = generateHash(currentMetric, key)

		metricsArr = append(metricsArr, currentMetric)
	}

	// add random value metric to arr
	currentMetric.ID = randomValueMetric.Name
	currentMetric.MType = randomValueMetric.Typename
	currentMetric.Value = &randomValueMetric.Value
	currentMetric.Hash = generateHash(currentMetric, key)
	metricsArr = append(metricsArr, currentMetric)

	// bad code - resetting for correct next assignments
	currentMetric.Value = nil

	// add counter metric to arr
	currentMetric.ID = pollCountMetric.Name
	currentMetric.MType = pollCountMetric.Typename
	currentMetric.Delta = &pollCountMetric.Value
	currentMetric.Hash = generateHash(currentMetric, key)
	metricsArr = append(metricsArr, currentMetric)

	return metricsArr
}

func generateHash(currentMetric common.Metrics, key string) string {
	hash, err := currentMetric.ComputeHash(key)
	if err != nil {
		return ""
	} else {
		return string(hash)
	}
}
