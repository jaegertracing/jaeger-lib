package xkit

import (
	kit "github.com/go-kit/kit/metrics"

	"github.com/uber/jaeger-lib/metrics"
)

// Factory provides a unified interface for creating named metrics
// from various go-kit metrics implementations.
type Factory interface {
	Counter(name string) kit.Counter
	Gauge(name string) kit.Gauge
	Histogram(name string) kit.Histogram
}

// Wrap is used to create an adapter from xkit.Factory to metrics.Factory.
func Wrap(namespace string, f Factory) metrics.Factory {
	return &factory{
		scope:   namespace,
		factory: f,
	}
}

type factory struct {
	scope   string
	tags    map[string]string
	factory Factory
}

func (f *factory) subScope(name string) string {
	if f.scope == "" {
		return name
	}
	if name == "" {
		return f.scope
	}
	return f.scope + "." + name
}

func (f *factory) Counter(name string, tags map[string]string) metrics.Counter {
	counter := f.factory.Counter(f.subScope(name))
	tagsList := f.tagsList(tags)
	if len(tagsList) > 0 {
		counter = counter.With(tagsList...)
	}
	return NewCounter(counter)
}

func (f *factory) Timer(name string, tags map[string]string) metrics.Timer {
	hist := f.factory.Histogram(f.subScope(name))
	tagsList := f.tagsList(tags)
	if len(tagsList) > 0 {
		hist = hist.With(tagsList...)
	}
	return NewTimer(hist)
}

func (f *factory) Gauge(name string, tags map[string]string) metrics.Gauge {
	gauge := f.factory.Gauge(f.subScope(name))
	tagsList := f.tagsList(tags)
	if len(tagsList) > 0 {
		gauge = gauge.With(tagsList...)
	}
	return NewGauge(gauge)
}

func (f *factory) Namespace(name string, tags map[string]string) metrics.Factory {
	return &factory{
		scope:   f.subScope(name),
		tags:    f.mergeTags(tags),
		factory: f.factory,
	}
}

func (f *factory) tagsList(a map[string]string) []string {
	m := f.mergeTags(a)
	ret := make([]string, 0, 2*len(m))
	for k, v := range m {
		ret = append(ret, k, v)
	}
	return ret
}

func (f *factory) mergeTags(tags map[string]string) map[string]string {
	ret := make(map[string]string, len(f.tags)+len(tags))
	for k, v := range f.tags {
		ret[k] = v
	}
	for k, v := range tags {
		ret[k] = v
	}
	return ret
}
