package influx

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/influx"

	"github.com/uber/jaeger-lib/metrics/go-kit"
)

// NewFactory creates a new metrics factory using go-kit influx package.
func NewFactory(client *influx.Influx) xkit.Factory {
	return factory{
		client: client,
	}
}

type factory struct {
	buckets int
	client *influx.Influx
}

func (f factory) Counter(name string) metrics.Counter {
	return f.client.NewCounter(name)
}

func (f factory) Histogram(name string) metrics.Histogram {
	return f.client.NewHistogram(name)
}

func (f factory) Gauge(name string) metrics.Gauge {
	return f.client.NewGauge(name)
}

func (f factory) Capabilities() xkit.Capabilities {
	return xkit.Capabilities{Tagging: true}
}
