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

func (d utilizationData) toCanonical() []model.Metric {
	var modelData []model.Metric

	modelData = append(modelData,
		d.TotalMemory,
		d.FreeMemory,
	)
	modelData = append(modelData, d.CPUUtilizations...)

	return modelData
}
