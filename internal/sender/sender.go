package sender

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"log"
	"net/http"
)

const (
	host = "127.0.0.1"
	port = "8080"
)

func SendGaugeMetric(client *http.Client, m metrics.GaugeMetric) error {
	url := fmt.Sprintf("http://%s:%s/update/%s/%s/%f", host, port, m.Typename, m.Name, m.Value)
	return Send(client, url, m.Base)
}

func SendCounterMetric(client *http.Client, m metrics.CounterMetric) error {
	url := fmt.Sprintf("http://%s:%s/update/%s/%s/%d", host, port, m.Typename, m.Name, m.Value)
	return Send(client, url, m.Base)
}

func Send(client *http.Client, url string, m metrics.Base) error {
	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		log.Printf("Failed to \"Post\" metric \"%s\" of type \"%s\"", m.Name, m.Typename)
		return err
	}
	defer resp.Body.Close()
	return err
}
