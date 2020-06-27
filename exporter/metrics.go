package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	GetName() string
	GetNamespace() string
	GetSubsystem() string
	Descs() (map[string]*prometheus.Desc, error)
	Metrics(ds map[string]*prometheus.Desc) ([]prometheus.Metric, error)
}
