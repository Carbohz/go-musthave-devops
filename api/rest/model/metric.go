package model

// представления для api
// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;
type GaugeMetricRequest struct {
	//address string
	//type string
	//name string
	//value float64/string
}

type GaugeMetricResponse struct {
	//
}


// converter from GaugeMetricRequest to model.GaugeMetric