package agent

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"log"
)

func (agent *Agent) sendMetrics() {
	go agent.sendMemStats()
	agent.sendGaugeMetric(agent.metrics.randomValue)
	go agent.sendCounterMetric(agent.metrics.pollCount)
}

func (agent *Agent) sendMemStats() {
	for _, m := range agent.metrics.memStats {
		go agent.sendGaugeMetric(m)
	}
}

func (agent *Agent) sendGaugeMetric(m model.GaugeMetric) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%f", agent.config.Address, m.Typename, m.Name, m.Value)
	return agent.sendMetric(url, m.Common)
}

func (agent *Agent) sendCounterMetric(m model.CounterMetric) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%d", agent.config.Address, m.Typename, m.Name, m.Value)
	return agent.sendMetric(url, m.Common)
}

func (agent *Agent) sendMetric(url string, m model.Common) error {
	resp, err := agent.client.Post(url, "text/plain", nil)
	if err != nil {
		log.Printf("Failed to \"Post\" request to update metric \"%s\" of type \"%s\"", m.Name, m.Typename)
		log.Printf("Error: %v", err)
		return err
	}
	defer resp.Body.Close()
	return err
}
