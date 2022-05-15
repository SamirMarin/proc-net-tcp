package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PromClient struct {
	Counter prometheus.Counter
}

func NewPromClient(name string, help string) *PromClient {
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})

	return &PromClient{
		Counter: counter,
	}
}

func (pc *PromClient) AddToCounter(value float64) {
	pc.Counter.Add(value)
}

func (pc *PromClient) IncToCounter() {
	pc.Counter.Inc()
}
