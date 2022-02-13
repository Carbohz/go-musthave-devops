package agent

import (
	"encoding/json"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/api/rest/models"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/go-resty/resty/v2"
	"log"
)

func (agent *Agent) sendMetrics() {
	go agent.sendMemStats()
	go agent.sendMetric(agent.metrics.randomValue)
	go agent.sendMetric(agent.metrics.pollCount)
}

func (agent *Agent) sendMetricsJSON() {
	go agent.sendMemStatsJSON()
	go agent.sendMetricJSON(agent.metrics.randomValue)
	go agent.sendMetricJSON(agent.metrics.pollCount)
}

func (agent *Agent) sendMemStats() {
	for _, m := range agent.metrics.memStats {
		go agent.sendMetric(m)
	}
}

func (agent *Agent) sendMemStatsJSON() {
	for _, m := range agent.metrics.memStats {
		go agent.sendMetricJSON(m)
	}
}

func (agent *Agent) sendMetric(m model.Metric) error {
	var url string

	if m.Delta.Present() {
		delta := m.MustGetInt()
		url = fmt.Sprintf("http://%s/update/%s/%s/%d", agent.config.Address, model.KCounter, m.Name, delta)
	} else {
		value := m.MustGetFloat()
		url = fmt.Sprintf("http://%s/update/%s/%s/%f", agent.config.Address, model.KCounter, m.Name, value)
	}

	_, err := agent.client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		log.Printf("Failed to \"Post\" request to update metric \"%s\" of type \"%s\"", m.Name, m.Type)
		log.Printf("Error: %v", err)
		return err
	}
	//logResponse(resp, err)

	return err
}


func (agent *Agent) sendMetricJSON(m model.Metric) error {
	url := fmt.Sprintf("http://%s/update/", agent.config.Address)

	metricToSend, err := models.FromModelMetrics(m)
	if err != nil {
		log.Printf("Error occured in agent.sendMetricJSON: %v", err)
		return fmt.Errorf("sendMetricJSON failed: %w", err)
	}

	//if agent.config.Key != "" {
	metricToSend.Hash = metricToSend.GenerateHash(agent.config.Key)
	//}

	rawJSON, err := json.Marshal(metricToSend)
	if err != nil {
		log.Printf("Error occured during metrics marshalling: %v", err)
	}
	log.Printf("Sending following body %v in JSON request", string(rawJSON))

	_, err = agent.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(rawJSON).
		EnableTrace().
		Post(url)
	if err != nil {
		log.Printf("Failed to \"Post\" json to update metric \"%s\" of type \"%s\"", m.Name, m.Type)
		log.Printf("Error: %v", err)
		return err
	}

	//logResponse(resp, err)

	return err
}

func logResponse(resp *resty.Response, err error) {
	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}