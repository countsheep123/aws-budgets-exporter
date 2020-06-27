package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Collector struct {
	up prometheus.Gauge

	ds map[string]*prometheus.Desc
	ms func(ds map[string]*prometheus.Desc) ([]prometheus.Metric, error)
}

func new(m Metrics) (*Collector, error) {
	ds, err := m.Descs()
	if err != nil {
		return nil, err
	}

	return &Collector{
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("%s_%s", m.GetNamespace(), m.GetSubsystem()),
			Name:      "up",
			Help:      "up",
		}),
		ds: ds,
		ms: m.Metrics,
	}, nil
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.up.Describe(ch)

	for _, desc := range c.ds {
		ch <- desc
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	zap.S().Info("collect")

	ms, err := c.ms(c.ds)
	if err != nil {
		c.up.Set(0)
		ch <- c.up
		zap.S().Error(err)
		return
	}

	c.up.Set(1)
	ch <- c.up

	for _, m := range ms {
		ch <- m
	}
}
