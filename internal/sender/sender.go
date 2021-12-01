package sender

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"net/http"
)

const (
	host = "127.0.0.1"
	port = "8080"
)

func Send(client *http.Client, m metrics.Metric) (*http.Response, error) {
	url := generateURL(m)
	resp, err := client.Post(url, "text/plain", nil)
	return resp, err
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