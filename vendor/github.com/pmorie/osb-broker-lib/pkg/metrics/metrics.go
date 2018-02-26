package metrics

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

const actionsMetricName = "osb_actions_total"

type OSBMetricsCollector struct {
	Actions *prom.CounterVec
}

func New() *OSBMetricsCollector {
	return &OSBMetricsCollector{
		Actions: prom.NewCounterVec(prom.CounterOpts{
			Name: actionsMetricName,
			Help: "Total amount of actions requested.",
		}, []string{"action"}),
	}
}

// Describe returns all descriptions of the collector.
func (c *OSBMetricsCollector) Describe(ch chan<- *prom.Desc) {
	c.Actions.Describe(ch)
}

// Collect returns the current state of all metrics of the collector.
func (c *OSBMetricsCollector) Collect(ch chan<- prom.Metric) {
	c.Actions.Collect(ch)
}
