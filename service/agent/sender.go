package agent

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"log"
)

func (agent *Agent) sendMetrics() {
	go agent.sendMemStats()
	go agent.sendMetric(agent.metrics.randomValue)
	go agent.sendMetric(agent.metrics.pollCount)
}

func (agent *Agent) sendMemStats() {
	for _, m := range agent.metrics.memStats {
		go agent.sendMetric(m)
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

	resp, err := agent.client.Post(url, "text/plain", nil)
	if err != nil {
		log.Printf("Failed to \"Post\" request to update metric \"%s\" of type \"%s\"", m.Name, m.Type)
		log.Printf("Error: %v", err)
		return err
	}
	defer resp.Body.Close()
	return err
}
