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
	Capabilities() Capabilities
}

// Capabilities describes capabilities of a specific metrics factory
type Capabilities struct {
	// Tagging indicates whether the factory has the capability for tagged metrics
	Tagging bool
}

// Wrap is used to create an adapter from xkit.Factory to metrics.Factory.
func Wrap(namespace string, f Factory) metrics.Factory {
	return &factory{
		scope:    namespace,
		factory:  f,
		scopeSep: ".",
		tagsSep:  ".",
		tagKVSep: "_",
	}
}

type factory struct {
	scope    string
	tags     map[string]string
	factory  Factory
	scopeSep string
	tagsSep  string
	tagKVSep string
}

func (f *factory) subScope(name string) string {
	if f.scope == "" {
		return name
	}
	if name == "" {
		return f.scope
	}
	return f.scope + f.scopeSep + name
}

// nameAndTagsList returns a name and tags list for the new metrics.
// The name is a concatenation of nom and the current factory scope.
// The tags list is a flattened list of passed tags merged with factory tags.
// If the underlying factory does not support tags, then the tags are
// transformed into a string and appended to the name.
func (f *factory) nameAndTagsList(nom string, tags map[string]string) (name string, tagsList []string) {
	mergedTags := f.mergeTags(tags)
	name = f.subScope(nom)
	tagsList = f.tagsList(mergedTags)
	if len(tagsList) == 0 {
		return
	}
	if f.factory.Capabilities().Tagging {
		return
	}
	name = metrics.GetKey(name, mergedTags, f.tagsSep, f.tagKVSep)
	tagsList = nil
	return
}

func (f *factory) Counter(name string, tags map[string]string) metrics.Counter {
	name, tagsList := f.nameAndTagsList(name, tags)
	counter := f.factory.Counter(name)
	if len(tagsList) > 0 {
		counter = counter.With(tagsList...)
	}
	return NewCounter(counter)
}

func (f *factory) Timer(name string, tags map[string]string) metrics.Timer {
	name, tagsList := f.nameAndTagsList(name, tags)
	hist := f.factory.Histogram(name)
	if len(tagsList) > 0 {
		hist = hist.With(tagsList...)
	}
	return NewTimer(hist)
}

func (f *factory) Gauge(name string, tags map[string]string) metrics.Gauge {
	name, tagsList := f.nameAndTagsList(name, tags)
	gauge := f.factory.Gauge(name)
	if len(tagsList) > 0 {
		gauge = gauge.With(tagsList...)
	}
	return NewGauge(gauge)
}

func (f *factory) Namespace(name string, tags map[string]string) metrics.Factory {
	return &factory{
		scope:    f.subScope(name),
		tags:     f.mergeTags(tags),
		factory:  f.factory,
		scopeSep: f.scopeSep,
		tagsSep:  f.tagsSep,
		tagKVSep: f.tagKVSep,
	}
}

func (f *factory) tagsList(a map[string]string) []string {
	ret := make([]string, 0, 2*len(a))
	for k, v := range a {
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
