package agent

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/markphelps/optional"
	"log"
	"net/http"
	"time"
)

type metrics struct {
	memStats    []model.Metric
	randomValue model.Metric
	pollCount   model.Metric
}

type Agent struct {
	config Config
	metrics metrics
	client *http.Client
}

func NewAgent(config Config) (*Agent, error) {
	client := http.Client{Timeout: 2 * time.Second}

	var m metrics
	pollCount := optional.NewInt64(0)
	m.pollCount = model.Metric{Name: "PollCount", Type: model.KCounter, Delta: pollCount}

	agent := &Agent{
		config: config,
		metrics: m,
		client: &client,
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
			//agent.sendMetrics()
			agent.sendMetricsJSON()
		case <-ctx.Done():
			return nil
		}
	}
}

