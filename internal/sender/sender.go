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

func Send(client *http.Client, m metrics.Metric) error {
	url := generateURL(m)
	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		log.Printf("Failed to Post metric \"%s\" of type \"%s\"", m.Name, m.Typename)
	}
	defer resp.Body.Close()
	return err
}

func generateURL(m metrics.Metric) string {
	if m.Typename == metrics.Gauge {
		return fmt.Sprintf("http://%s:%s/update/%s/%s/%f", host, port, m.Typename, m.Name, m.Value)
	} else if m.Typename == metrics.Counter {
		return fmt.Sprintf("http://%s:%s/update/%s/%s/%d", host, port, m.Typename, m.Name, int64(m.Value))
	} else {
		return ""
	}
}