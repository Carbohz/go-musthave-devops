package server

import "github.com/Carbohz/go-musthave-devops/internal/metrics"

type Instance struct {
	Cfg Config
	is  *internalStorage // `is` is shortcut for internal storage
	dbs *dbStorage       // `dbs` is shortcut for database storage
}

// insertion into storage
func (instance Instance) StoreGaugeMetric(name string, value float64) {
	if instance.Cfg.DBPath != "" {
		// working with db
	} else {
		// working with internal storage
		instance.is.StoreGaugeMetric(name, value)
	}
}

func (instance Instance) StoreCounterMetric(name string, value int64) {
	if instance.Cfg.DBPath != "" {
		// working with db
	} else {
		// working with internal storage
		instance.is.StoreCounterMetric(name, value)
	}
}

// storage lookup
func (instance Instance) FindGaugeMetric(name string) (float64, error) {
	if instance.Cfg.DBPath != "" {
		// working with db
		return -1.0, nil
	} else {
		// working with internal storage
		return instance.is.FindGaugeMetric(name)
	}
}

func (instance Instance) FindCounterMetric(name string) (int64, error) {
	if instance.Cfg.DBPath != "" {
		// working with db
		return -1, nil
	} else {
		// working with internal storage
		return instance.is.FindCounterMetric(name)
	}
}

// load data
func (instance Instance) LoadGaugeMetrics() map[string]metrics.GaugeMetric {
	if instance.Cfg.DBPath != "" {
		// working with db
		return make(map[string]metrics.GaugeMetric)
	} else {
		// working with internal storage
		return instance.is.LoadGaugeMetrics()
	}
}

func (instance Instance) LoadCounterMetrics() map[string]metrics.CounterMetric {
	if instance.Cfg.DBPath != "" {
		// working with db
		return make(map[string]metrics.CounterMetric)
	} else {
		// working with internal storage
		return instance.is.LoadCounterMetrics()
	}
}

// update data
func (instance Instance) UpdateGaugeMetrics(name string, value float64) {
	if instance.Cfg.DBPath != "" {
		// working with db

	} else {
		// working with internal storage
		instance.is.UpdateGaugeMetrics(name, value)
	}
}

func (instance Instance) UpdateCounterMetrics(name string, value int64) {
	if instance.Cfg.DBPath != "" {
		// working with db

	} else {
		// working with internal storage
		instance.is.UpdateCounterMetrics(name, value)
	}
}