package metrics

import (
	"os"

	"github.com/caarlos0/env"
	"github.com/countsheep123/aws-budgets-exporter/aws/budgets"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type config struct {
	ConfigPath      string `env:"CONFIG_PATH" envDefault:"/opt/config.yaml"`
	RoleSessionName string `env:"ROLE_SESSION_NAME,required"`
}

type account struct {
	ID     string            `yaml:"id,omitempty"`
	Labels map[string]string `yaml:"labels,omitempty"`
}

type request struct {
	RoleArn  string     `yaml:"role_arn,omitempty"`
	Accounts []*account `yaml:"accounts,omitempty"`
}

func New() (*Metrics, error) {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	file, err := os.Open(cfg.ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var requests []*request
	if err := yaml.NewDecoder(file).Decode(&requests); err != nil {
		return nil, err
	}

	targets := []*target{}
	labelMap := map[string]bool{}

	for _, r := range requests {
		cli, err := budgets.New(r.RoleArn, cfg.RoleSessionName)
		if err != nil {
			return nil, err
		}
		for _, a := range r.Accounts {
			targets = append(targets, &target{
				cli:       cli,
				accountID: a.ID,
				labels:    a.Labels,
			})
			for k, _ := range a.Labels {
				labelMap[k] = true
			}
		}
	}

	labels := []string{}
	for k, _ := range labelMap {
		labels = append(labels, k)
	}

	return &Metrics{
		name:      "aws_budgets_exporter",
		namespace: "aws",
		subsystem: "budgets",
		labels:    labels,
		targets:   targets,
	}, nil
}

type target struct {
	cli       *budgets.Client
	accountID string
	labels    map[string]string
}

type Metrics struct {
	name      string
	namespace string
	subsystem string

	labels  []string
	targets []*target
}

func (m *Metrics) GetName() string {
	return m.name
}

func (m *Metrics) GetNamespace() string {
	return m.namespace
}

func (m *Metrics) GetSubsystem() string {
	return m.subsystem
}

func (m *Metrics) Descs() (map[string]*prometheus.Desc, error) {
	descs := map[string]*prometheus.Desc{}

	{
		name := "actual_spend"
		labels := []string{"budget_name", "account_id"}
		labels = append(labels, m.labels...)
		descs[name] = prometheus.NewDesc(
			prometheus.BuildFQName(m.namespace, m.subsystem, name),
			name,
			labels,
			nil,
		)
	}
	{
		name := "forecasted_spend"
		labels := []string{"budget_name", "account_id"}
		labels = append(labels, m.labels...)
		descs[name] = prometheus.NewDesc(
			prometheus.BuildFQName(m.namespace, m.subsystem, name),
			name,
			labels,
			nil,
		)
	}

	return descs, nil
}

func (m *Metrics) Metrics(ds map[string]*prometheus.Desc) ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	for _, t := range m.targets {
		budgets, err := t.cli.GetBudgets(t.accountID)
		if err != nil {
			return nil, err
		}

		ls := []string{}
		for _, l := range m.labels {
			ls = append(ls, t.labels[l])
		}

		for _, budget := range budgets {
			labels := []string{budget.Name, t.accountID}
			labels = append(labels, ls...)

			if budget.ActualSpend != nil {
				metrics = append(metrics, prometheus.MustNewConstMetric(
					ds["actual_spend"],
					prometheus.GaugeValue,
					*budget.ActualSpend,
					labels...,
				))

				zap.S().Info("actual_spend",
					"value", *budget.ActualSpend,
					"labels", labels,
				)
			} else {
				zap.S().Info("actual_spend",
					"value", "NODATA",
					"labels", labels,
				)
			}

			if budget.ForecastedSpend != nil {
				metrics = append(metrics, prometheus.MustNewConstMetric(
					ds["forecasted_spend"],
					prometheus.GaugeValue,
					*budget.ForecastedSpend,
					labels...,
				))

				zap.S().Info("forecasted_spend",
					"value", *budget.ForecastedSpend,
					"labels", labels,
				)
			} else {
				zap.S().Info("forecasted_spend",
					"value", "NODATA",
					"labels", labels,
				)
			}
		}
	}

	return metrics, nil
}
