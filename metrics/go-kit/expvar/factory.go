package expvar

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"

	"github.com/uber/jaeger-lib/metrics/go-kit"
)

// NewFactory creates a new metrics factory using go-kit expvar package.
// buckets is the number of buckets to be used in histograms.
func NewFactory(buckets int) xkit.Factory {
	return factory{
		buckets: buckets,
	}
}

type factory struct {
	buckets int
}

func (f factory) Counter(name string) metrics.Counter {
	return expvar.NewCounter(name)
}

func (f factory) Histogram(name string) metrics.Histogram {
	return expvar.NewHistogram(name, f.buckets)
}

func (f factory) Gauge(name string) metrics.Gauge {
	return expvar.NewGauge(name)
}

func (f factory) Capabilities() xkit.Capabilities {
	return xkit.Capabilities{Tagging: false}
}
