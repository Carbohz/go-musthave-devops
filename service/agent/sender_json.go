package agent

import (
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/Carbohz/go-musthave-devops/model"
	"log"
)

func (a *Agent) sendMetricsWithJSON() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	go a.sendMetricsSliceWithJSON(a.metrics.memStats)
	go a.sendSingleMetricWithJSON(a.metrics.randomValue)
	go a.sendSingleMetricWithJSON(a.metrics.pollCount)
	go a.sendMetricsSliceWithJSON(toModelUtilizationData(a.metrics.utilization))
}

func (a *Agent) sendMetricsBatchWithJSON() error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var metricsArr []models.Metrics

	for _, m := range a.metrics.memStats {
		v, _ := models.NewMetricFromCanonical(m)
		v.Hash = v.GenerateHash(a.config.Key)
		metricsArr = append(metricsArr, v)
	}

	randomValue, _ := models.NewMetricFromCanonical(a.metrics.randomValue)
	randomValue.Hash = randomValue.GenerateHash(a.config.Key)
	metricsArr = append(metricsArr, randomValue)

	pollCount, _ := models.NewMetricFromCanonical(a.metrics.pollCount)
	pollCount.Hash = pollCount.GenerateHash(a.config.Key)
	metricsArr = append(metricsArr, pollCount)

	utilMetricsArr := toModelUtilizationData(a.metrics.utilization)
	for _, m := range utilMetricsArr {
		v, _ := models.NewMetricFromCanonical(m)
		v.Hash = v.GenerateHash(a.config.Key)
		metricsArr = append(metricsArr, v)
	}

	rawJSON, err := json.Marshal(metricsArr)
	if err != nil {
		log.Printf("Error occured during metrics marshalling: %v", err)
	}
	log.Printf("Sending following body %v in JSON request", string(rawJSON))

	url := fmt.Sprintf("http://%s/updates/", a.config.Address)

	_, err = a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(rawJSON).
		EnableTrace().
		Post(url)
	if err != nil {
		log.Println("Failed to send batch of metrics")
		log.Printf("Error: %v", err)
		return err
	}

	return nil
}

func (a *Agent) sendSingleMetricWithJSON(m model.Metric) error {
	url := fmt.Sprintf("http://%s/update/", a.config.Address)

	metricToSend, err := models.NewMetricFromCanonical(m)
	if err != nil {
		log.Printf("Error occured in a.sendSingleMetricWithJSON: %v", err)
		return fmt.Errorf("sendSingleMetricWithJSON failed: %w", err)
	}

	metricToSend.Hash = metricToSend.GenerateHash(a.config.Key)

	rawJSON, err := json.Marshal(metricToSend)
	if err != nil {
		log.Printf("Error occured during metrics marshalling: %v", err)
	}
	log.Printf("Sending following body %v in JSON request", string(rawJSON))

	_, err = a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(rawJSON).
		EnableTrace().
		Post(url)
	if err != nil {
		log.Printf("Failed to \"Post\" json to update metric \"%s\" of type \"%s\"", m.Name, m.Type)
		log.Printf("Error: %v", err)
		return err
	}

	return err
}

func (a *Agent) sendMetricsSliceWithJSON(slice []model.Metric) {
	for _, m := range slice {
		go a.sendSingleMetricWithJSON(m)
	}

}
