package agent

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"log"
)

func (a *Agent) sendMetricsWithURL() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	go a.sendMetricsSliceWithURL(a.metrics.memStats)
	go a.sendSingleMetricWithURL(a.metrics.randomValue)
	go a.sendSingleMetricWithURL(a.metrics.pollCount)
	go a.sendMetricsSliceWithURL(toModelUtilizationData(a.metrics.utilization))
}

func (a *Agent) sendSingleMetricWithURL(m model.Metric) error {
	var url string

	if m.Delta.Present() {
		delta := m.MustGetInt()
		url = fmt.Sprintf("http://%s/update/%s/%s/%d", a.config.Address, model.KCounter, m.Name, delta)
	} else {
		value := m.MustGetFloat()
		url = fmt.Sprintf("http://%s/update/%s/%s/%.6f", a.config.Address, model.KGauge, m.Name, value)
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

func (a *Agent) sendMetricsSliceWithURL(slice []model.Metric ) {
	for _, m := range slice {
		go a.sendSingleMetricWithURL(m)
	}
}


