// Copyright (c) 2017 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/uber/jaeger-lib/metrics"
)

// Factory implements metrics.Factory backed my Prometheus registry.
type Factory struct {
	scope string
	tags  map[string]string
	cache *vectorCache
}

// New creates a Factory backed by Prometheus registry.
// Typically the first argument should be prometheus.DefaultRegisterer.
func New(registerer prometheus.Registerer) *Factory {
	return newFactory(newVectorCache(registerer), "", nil)
}

func newFactory(cache *vectorCache, scope string, tags map[string]string) *Factory {
	return &Factory{
		cache: cache,
		scope: scope,
		tags:  tags,
	}
}

// Counter implements Counter of metrics.Factory.
func (f *Factory) Counter(name string, tags map[string]string) metrics.Counter {
	name = f.subScope(name)
	tags = f.mergeTags(tags)
	opts := prometheus.CounterOpts{
		Name: name,
		Help: name,
	}
	labelNames := f.tagNames(tags)
	cv := f.cache.getOrMakeCounterVec(opts, labelNames)
	return &counter{
		counter: cv.WithLabelValues(f.tagsAsLabelValues(labelNames, tags)...),
	}
}

// Gauge implements Gauge of metrics.Factory.
func (f *Factory) Gauge(name string, tags map[string]string) metrics.Gauge {
	name = f.subScope(name)
	tags = f.mergeTags(tags)
	opts := prometheus.GaugeOpts{
		Name: name,
		Help: name,
	}
	labelNames := f.tagNames(tags)
	gv := f.cache.getOrMakeGaugeVec(opts, labelNames)
	return &gauge{
		gauge: gv.WithLabelValues(f.tagsAsLabelValues(labelNames, tags)...),
	}
}

// Timer implements Timer of metrics.Factory.
func (f *Factory) Timer(name string, tags map[string]string) metrics.Timer {
	name = f.subScope(name)
	tags = f.mergeTags(tags)
	opts := prometheus.HistogramOpts{
		Name: name,
		Help: name,
	}
	labelNames := f.tagNames(tags)
	hv := f.cache.getOrMakeHistogramVec(opts, labelNames)
	return &timer{
		histogram: hv.WithLabelValues(f.tagsAsLabelValues(labelNames, tags)...),
	}
}

// Namespace implements Namespace of metrics.Factory.
func (f *Factory) Namespace(name string, tags map[string]string) metrics.Factory {
	return newFactory(f.cache, f.subScope(name), f.mergeTags(tags))
}

type counter struct {
	counter prometheus.Counter
}

func (c *counter) Inc(v int64) {
	c.counter.Add(float64(v))
}

type gauge struct {
	gauge prometheus.Gauge
}

func (g *gauge) Update(v int64) {
	g.gauge.Set(float64(v))
}

type timer struct {
	histogram prometheus.Histogram
}

func (t *timer) Record(v time.Duration) {
	t.histogram.Observe(float64(v.Nanoseconds()) / float64(time.Second/time.Nanosecond))
}

func (f *Factory) subScope(name string) string {
	if f.scope == "" {
		return name
	}
	if name == "" {
		return f.scope
	}
	return f.scope + ":" + name
}

func (f *Factory) mergeTags(tags map[string]string) map[string]string {
	ret := make(map[string]string, len(f.tags)+len(tags))
	for k, v := range f.tags {
		ret[k] = v
	}
	for k, v := range tags {
		ret[k] = v
	}
	return ret
}

func (f *Factory) tagNames(tags map[string]string) []string {
	ret := make([]string, 0, len(tags))
	for k := range tags {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

func (f *Factory) tagsAsLabelValues(labels []string, tags map[string]string) []string {
	ret := make([]string, 0, len(tags))
	for _, l := range labels {
		ret = append(ret, tags[l])
	}
	return ret
}
