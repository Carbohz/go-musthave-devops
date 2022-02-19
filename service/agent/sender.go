package agent

import (
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/Carbohz/go-musthave-devops/model"
	"log"
)

func (agent *Agent) sendMetrics() {
	go agent.sendMemStats()
	go agent.sendMetric(agent.metrics.randomValue)
	go agent.sendMetric(agent.metrics.pollCount)
}

func (agent *Agent) sendMetricsJSON() {
	go agent.sendMemStatsJSON()
	go agent.sendMetricJSON(agent.metrics.randomValue)
	go agent.sendMetricJSON(agent.metrics.pollCount)
}

func (agent *Agent) sendMetricsBatch() error {
	var metricsArr []models.Metrics

	for _, m := range agent.metrics.memStats {
		v, _ := models.NewMetricFromCanonical(m)
		v.Hash = v.GenerateHash(agent.config.Key)
		metricsArr = append(metricsArr, v)
	}

	randomValue, _ := models.NewMetricFromCanonical(agent.metrics.randomValue)
	randomValue.Hash = randomValue.GenerateHash(agent.config.Key)
	metricsArr = append(metricsArr, randomValue)

	pollCount, _ := models.NewMetricFromCanonical(agent.metrics.pollCount)
	pollCount.Hash = pollCount.GenerateHash(agent.config.Key)
	metricsArr = append(metricsArr, pollCount)

	rawJSON, err := json.Marshal(metricsArr)
	if err != nil {
		log.Printf("Error occured during metrics marshalling: %v", err)
	}
	log.Printf("Sending following body %v in JSON request", string(rawJSON))

	url := fmt.Sprintf("http://%s/updates/", agent.config.Address)

	_, err = agent.client.R().
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

func (agent *Agent) sendMemStats() {
	for _, m := range agent.metrics.memStats {
		go agent.sendMetric(m)
	}
}

func (agent *Agent) sendMemStatsJSON() {
	for _, m := range agent.metrics.memStats {
		go agent.sendMetricJSON(m)
	}
}

func (agent *Agent) sendMetric(m model.Metric) error {
	var url string

	if m.Delta.Present() {
		delta := m.MustGetInt()
		url = fmt.Sprintf("http://%s/update/%s/%s/%d", agent.config.Address, model.KCounter, m.Name, delta)
	} else {
		value := m.MustGetFloat()
		url = fmt.Sprintf("http://%s/update/%s/%s/%f", agent.config.Address, model.KCounter, m.Name, value)
	}

	_, err := agent.client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		log.Printf("Failed to \"Post\" request to update metric \"%s\" of type \"%s\"", m.Name, m.Type)
		log.Printf("Error: %v", err)
		return err
	}

	return err
}


func (agent *Agent) sendMetricJSON(m model.Metric) error {
	url := fmt.Sprintf("http://%s/update/", agent.config.Address)

	metricToSend, err := models.NewMetricFromCanonical(m)
	if err != nil {
		log.Printf("Error occured in agent.sendMetricJSON: %v", err)
		return fmt.Errorf("sendMetricJSON failed: %w", err)
	}

	metricToSend.Hash = metricToSend.GenerateHash(agent.config.Key)

	rawJSON, err := json.Marshal(metricToSend)
	if err != nil {
		log.Printf("Error occured during metrics marshalling: %v", err)
	}
	log.Printf("Sending following body %v in JSON request", string(rawJSON))

	_, err = agent.client.R().
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
