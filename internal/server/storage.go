package server

import (
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"log"
	"os"
)

type internalStorage struct {
	GaugeMetrics   map[string]metrics.GaugeMetric
	CounterMetrics map[string]metrics.CounterMetric
}

// insertion into storage
func (s internalStorage) StoreGaugeMetric(name string, value float64) {
	s.GaugeMetrics[name] = metrics.GaugeMetric{
		Base:  metrics.Base{Name: name, Typename: metrics.Gauge},
		Value: value}
}

func (s internalStorage) StoreCounterMetric(name string, value int64) {
	s.CounterMetrics[name] = metrics.CounterMetric{
		Base:  metrics.Base{Name: name, Typename: metrics.Counter},
		Value: s.CounterMetrics[name].Value + value}
}

// storage lookup
func (s internalStorage) FindGaugeMetric(name string) (float64, error) {
	if value, found := s.GaugeMetrics[name]; found {
		return value.Value, nil
	}
	err := fmt.Errorf("Unknown metric \"%s\" of type \"gauge\"", name)
	return -1.0, err
}

func (s internalStorage) FindCounterMetric(name string) (int64, error) {
	if value, found := s.CounterMetrics[name]; found {
		return value.Value, nil
	}
	err := fmt.Errorf("Unknown metric \"%s\" of type \"counter\"", name)
	return -1, err
}

// get data
func (s internalStorage) GetGaugeMetrics() map[string]metrics.GaugeMetric {
	return s.GaugeMetrics
}

func (s internalStorage) GetCounterMetrics() map[string]metrics.CounterMetric {
	return s.CounterMetrics
}

// update data
func (s internalStorage) UpdateGaugeMetrics(name string, value float64) {
	s.GaugeMetrics[name] = metrics.GaugeMetric{
		Base:  metrics.Base{Name: name, Typename: metrics.Gauge},
		Value: value}
}

func (s internalStorage) UpdateCounterMetrics(name string, value int64) {
	s.CounterMetrics[name] = metrics.CounterMetric{
		Base:  metrics.Base{Name: name, Typename: metrics.Gauge},
		Value: s.CounterMetrics[name].Value + value}
}

// dump data to file
func (s internalStorage) DumpMetricsToFile(instance Instance) {
	cfg := instance.Cfg

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(cfg.StoreFile, flag, 0644)
	if err != nil {
		log.Fatal("Can't open file for dumping: ", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	//internalStorage := InternalStorage{
	//	GaugeMetrics:   gaugeMetricsStorage,
	//	CounterMetrics: counterMetricsStorage,
	//}

	if err := encoder.Encode(s); err != nil {
		log.Fatal("Can't encode server's metrics: ", err)
	}
}

// load data from file
func (s internalStorage) LoadMetricsFromFile(instance Instance) {
	cfg := instance.Cfg
	log.Printf("Loading metrics from file %s", cfg.StoreFile)

	flag := os.O_RDONLY

	f, err := os.OpenFile(cfg.StoreFile, flag, 0)
	if err != nil {
		log.Print("Can't open file for loading metrics: ", err)
		return
	}
	defer f.Close()

	var internalStorage internalStorage

	if err := json.NewDecoder(f).Decode(&internalStorage); err != nil {
		log.Fatal("Can't decode metrics: ", err)
	}

	s.GaugeMetrics = internalStorage.GaugeMetrics
	s.CounterMetrics = internalStorage.CounterMetrics
	log.Printf("Metrics successfully loaded from file %s", cfg.StoreFile)
}
