package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
func (m Metrics) computeHash(key string) ([]byte, error) {
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

func (m Metrics) GenerateHash(key string) string {
	hash, err := m.computeHash(key)
	if err != nil {
		return ""
	} else {
		//1
		//return string(hash)

		// 2
		//return fmt.Sprintf("%x", hash)

		// 3
		return hex.EncodeToString(hash)
	}
}

func (m Metrics) CheckHash(key string) error {
	if key == "" {
		return nil
	}
	h, err := m.computeHash(key)
	if err != nil {
		return err
	}
	hashStr := hex.EncodeToString(h)
	if m.Hash != hashStr {
		return fmt.Errorf("hash value incorrect")
	}
	return nil
}
