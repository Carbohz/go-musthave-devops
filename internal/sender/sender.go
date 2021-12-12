package sender

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"log"
	"net/http"
)

func SendGaugeMetric(client *http.Client, m metrics.GaugeMetric, address string) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%f", address, m.Typename, m.Name, m.Value)
	return Send(client, url, m.Base)
}

func SendCounterMetric(client *http.Client, m metrics.CounterMetric, address string) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%d", address, m.Typename, m.Name, m.Value)
	return Send(client, url, m.Base)
}

func Send(client *http.Client, url string, m metrics.Base) error {
	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		log.Printf("Failed to \"Post\" request to update metric \"%s\" of type \"%s\"", m.Name, m.Typename)
		return err
	}
	defer resp.Body.Close()
	return err
}
