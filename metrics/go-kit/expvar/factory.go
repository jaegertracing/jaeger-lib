package expvar

import (
	kit "github.com/go-kit/kit/metrics/expvar"

	"github.com/uber/jaeger-lib/metrics"
	xkit "github.com/uber/jaeger-lib/metrics/go-kit"
)

// NewFactory creates a new metrics factory using go-kit expvar package.
// scope is the name that will be prepended to the names of all metrics
// created by this factory. buckets is the number of buckets to be used
// in histograms representing timers.
func NewFactory(scope string, buckets int) metrics.Factory {
	return factory{
		scope:   scope,
		buckets: buckets,
	}
}

type factory struct {
	scope   string
	buckets int
}

func (f factory) subScope(name string) string {
	return f.scope + "_" + name
}

func (f factory) Counter(name string, tags map[string]string) metrics.Counter {
	return xkit.NewCounter(kit.NewCounter(f.subScope(name)))
}

func (f factory) Timer(name string, tags map[string]string) metrics.Timer {
	return xkit.NewTimer(kit.NewHistogram(f.subScope(name), f.buckets))
}

func (f factory) Gauge(name string, tags map[string]string) metrics.Gauge {
	return xkit.NewGauge(kit.NewGauge(f.subScope(name)))
}

// Namespace is a no-op for expvar, since i
func (f factory) Namespace(name string, tags map[string]string) metrics.Factory {
	return factory{
		scope:   f.subScope(name),
		buckets: f.buckets,
	}
}
