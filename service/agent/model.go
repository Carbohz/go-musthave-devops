package agent

import (
	"github.com/Carbohz/go-musthave-devops/model"
	"sync"
	"time"
)

type utilizationData struct {
	mu              sync.Mutex
	TotalMemory     model.Metric
	FreeMemory      model.Metric
	CPUutilizations []model.Metric
	CPUtime         []float64
	CPUutilLastTime time.Time
}

type metrics struct {
	memStats    []model.Metric
	randomValue model.Metric
	pollCount   model.Metric
	utilization *utilizationData
}

func toModelUtilizationData(utilData *utilizationData) []model.Metric {
	var modelData []model.Metric

	modelData = append(modelData,
		utilData.TotalMemory,
		utilData.FreeMemory,
	)

	for _, d := range utilData.CPUutilizations {
		modelData = append(modelData, d)
	}


	//for i, t := range utilData.CPUtime {
	//	cpuTime := model.NewGaugeMetric("CPUtime" + strconv.Itoa(i+1), t)
	//	modelData = append(modelData, cpuTime)
	//}

	return modelData
}