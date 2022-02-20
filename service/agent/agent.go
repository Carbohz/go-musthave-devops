package agent

import (
	"context"
	configagent "github.com/Carbohz/go-musthave-devops/config/agent"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/markphelps/optional"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

//type utilizationData struct {
//	mu              sync.Mutex
//	TotalMemory     model.Metric
//	FreeMemory      model.Metric
//	CPUutilizations []model.Metric
//	CPUtime         []float64
//	CPUutilLastTime time.Time
//}
//
//type metrics struct {
//	memStats    []model.Metric
//	randomValue model.Metric
//	pollCount   model.Metric
//	utilization utilizationData
//}

type Agent struct {
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
	reportTicker := time.NewTicker(a.config.ReportInterval)

	for {
		select {
		case <-pollTicker.C:
			log.Println("Collecting Metrics")
			a.collectMetrics()
		case <-reportTicker.C:
			log.Println("Sending Metrics")
			go a.sendMetrics()
			//go a.sendMetricsJSON()
			//go a.sendMetricsBatch()
		case <-ctx.Done():
			return nil
		}
	}
}

