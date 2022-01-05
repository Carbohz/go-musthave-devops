package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

const (
	KGauge   = "gauge"
	KCounter = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

// ComputeHash calculates hash for metrics
func (m Metrics) ComputeHash(key string) ([]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("no key")
	}

	if m.ID == "" {
		return nil, fmt.Errorf("empty ID field")
	}

	toHash := ""

	if m.MType == KGauge {
		if m.Value == nil {
			return nil, fmt.Errorf("no value")
		}
		toHash = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	}

	if m.MType == KCounter {
		if m.Delta == nil {
			return nil, fmt.Errorf("no delta")
		}
		toHash = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}

	//h := hmac.New(sha256.New, []byte(key))
	//h.Write([]byte(toHash))
	//hash := h.Sum(nil)
	//return hash, nil

	//h := sha256.New()
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(toHash))
	hash := h.Sum(nil)
	//log.Printf("%x", hash)
	return hash, nil
}
