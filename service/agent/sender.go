package agent

import (
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/Carbohz/go-musthave-devops/model"
	"log"
)

func (a *Agent) sendMetrics() {
	//go a.sendMemStats()
	go a.sendMetricsSlice(a.metrics.memStats)

	go a.sendMetric(a.metrics.randomValue)
	go a.sendMetric(a.metrics.pollCount)

	//go a.sendMetric(a.metrics.utilization)
	go a.sendMetricsSlice(toModelUtilizationData(a.metrics.utilization))
}

func (a *Agent) sendMetricsJSON() {
	//go a.sendMemStatsJSON()
	go a.sendMetricsSliceJSON(a.metrics.memStats)

	go a.sendMetricJSON(a.metrics.randomValue)
	go a.sendMetricJSON(a.metrics.pollCount)

	go a.sendMetricsSliceJSON(toModelUtilizationData(a.metrics.utilization))
}

func (a *Agent) sendMetricsBatch() error {
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

func (a *Agent) sendMetricsSlice(slice []model.Metric ) {
	for _, m := range slice {
		go a.sendMetric(m)
	}
}

func (a *Agent) sendMemStats() {
	for _, m := range a.metrics.memStats {
		go a.sendMetric(m)
	}
}

func (a *Agent) sendMemStatsJSON() {
	for _, m := range a.metrics.memStats {
		go a.sendMetricJSON(m)
	}
}

func (a *Agent) sendMetricsSliceJSON(slice []model.Metric) {
	for _, m := range slice {
		go a.sendMetricJSON(m)
	}

}

//func (agent *Agent) sendCPUutilization() {
//	agent.metrics.utilization
//}

func (a *Agent) sendMetric(m model.Metric) error {
	var url string

	if m.Delta.Present() {
		delta := m.MustGetInt()
		url = fmt.Sprintf("http://%s/update/%s/%s/%d", a.config.Address, model.KCounter, m.Name, delta)
	} else {
		value := m.MustGetFloat()
		url = fmt.Sprintf("http://%s/update/%s/%s/%.20f", a.config.Address, model.KGauge, m.Name, value)
	}

	_, err := a.client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		log.Printf("Failed to \"Post\" request to update metric \"%s\" of type \"%s\"", m.Name, m.Type)
		log.Printf("Error: %v", err)
		return err
	}

	return err
}


func (a *Agent) sendMetricJSON(m model.Metric) error {
	url := fmt.Sprintf("http://%s/update/", a.config.Address)

	metricToSend, err := models.NewMetricFromCanonical(m)
	if err != nil {
		log.Printf("Error occured in a.sendMetricJSON: %v", err)
		return fmt.Errorf("sendMetricJSON failed: %w", err)
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
