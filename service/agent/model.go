package agent

import (
	"github.com/Carbohz/go-musthave-devops/model"
	"sync"
)

type utilizationData struct {
	mu              sync.Mutex
	TotalMemory     model.Metric
	FreeMemory      model.Metric
	CPUUtilizations []model.Metric
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
	modelData = append(modelData, utilData.CPUUtilizations...)

	return modelData
}
