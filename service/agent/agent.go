package agent

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/model"
	"log"
	"net/http"
	"time"
)

type metrics struct {
	memStats    []model.GaugeMetric
	randomValue model.GaugeMetric
	pollCount   model.CounterMetric
}

type Agent struct {
	config Config
	metrics metrics
	client http.Client
}

func NewAgent() (*Agent, error) {
	config := createConfig()
	client := http.Client{Timeout: 2 * time.Second}

	var m metrics
	m.pollCount = model.CounterMetric{Common: model.Common{Name: "PollCount", Typename: model.Counter}}

	agent := &Agent{
		config: config,
		metrics: m,
		client: client,
	}

	return agent, nil
}

func (agent *Agent) Run(ctx context.Context) error {
	pollTicker := time.NewTicker(agent.config.PollInterval)
	reportTicker := time.NewTicker(agent.config.ReportInterval)

	for {
		select {
		case <-pollTicker.C:
			log.Println("Collecting Metrics")
			agent.collectMetrics()
		case <-reportTicker.C:
			log.Println("Sending Metrics")
			agent.sendMetrics()
		case <-ctx.Done():
			return nil
		}
	}
}

