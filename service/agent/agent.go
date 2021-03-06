package agent

import (
	"context"
	configagent "github.com/Carbohz/go-musthave-devops/config/agent"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/markphelps/optional"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Agent struct {
	mu sync.RWMutex
	config configagent.AgentConfig
	metrics metrics
	client *resty.Client
}

func NewAgent(config configagent.AgentConfig) (*Agent, error) {
	var m metrics
	pollCount := optional.NewInt64(0)
	m.pollCount = model.Metric{Name: "PollCount", Type: model.KCounter, Delta: pollCount}

	client := resty.New()

	agent := &Agent{
		config: config,
		metrics: m,
		client: client,
	}

	return agent, nil
}

func (a *Agent) Run(ctx context.Context) error {
	pollTicker := time.NewTicker(a.config.PollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(a.config.ReportInterval)
	defer reportTicker.Stop()

	wg := sync.WaitGroup{}

	for {
		select {
		case <-pollTicker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				a.collectMetrics()
			}()
		case <-reportTicker.C:
			wg.Wait()
			go func() {
				a.sendMetricsBatchWithJSON()
			}()
		case <-ctx.Done():
			return nil
		}
	}
}

