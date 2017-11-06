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
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/uber/jaeger-lib/metrics"
)

// Factory implements metrics.Factory backed my Prometheus registry.
//
// See New.
type Factory struct {
	registerer prometheus.Registerer
	scope      string
	tags       map[string]string
	cVecs      map[string]*prometheus.CounterVec
	gVecs      map[string]*prometheus.GaugeVec
	hVecs      map[string]*prometheus.HistogramVec
	lock       sync.Mutex
}

// New creates a metrics.Factory backed by Prometheus registry.
// Typically the first argument should be prometheus.DefaultRegisterer.
func New(registerer prometheus.Registerer) *Factory {
	return newFactory(registerer, "", nil)
}

func newFactory(registerer prometheus.Registerer, scope string, tags map[string]string) *Factory {
	return &Factory{
		registerer: registerer,
		scope:      scope,
		tags:       tags,
		cVecs:      make(map[string]*prometheus.CounterVec),
		gVecs:      make(map[string]*prometheus.GaugeVec),
		hVecs:      make(map[string]*prometheus.HistogramVec),
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

	f.lock.Lock()
	defer f.lock.Unlock()

	cacheKey := strings.Join(append([]string{name}, labelNames...), "||")
	cv, cvExists := f.cVecs[cacheKey]
	if !cvExists {
		cv = prometheus.NewCounterVec(opts, labelNames)
		f.registerer.MustRegister(cv)
		f.cVecs[cacheKey] = cv
	}
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

	f.lock.Lock()
	defer f.lock.Unlock()

	cacheKey := strings.Join(append([]string{name}, labelNames...), "||")
	gv, gvExists := f.gVecs[cacheKey]
	if !gvExists {
		gv = prometheus.NewGaugeVec(opts, labelNames)
		f.registerer.MustRegister(gv)
		f.gVecs[cacheKey] = gv
	}
	return &gauge{
		gauge: gv.WithLabelValues(f.tagsAsLabelValues(labelNames, tags)...),
	}
}

// Timer implements Timer of metrics.Factory.
func (f *Factory) Timer(name string, tags map[string]string) metrics.Timer {
	panic("Timer() not implemented")
}

// Namespace implements Namespace of metrics.Factory.
func (f *Factory) Namespace(name string, tags map[string]string) metrics.Factory {
	return newFactory(f.registerer, f.subScope(name), f.mergeTags(tags))
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
		if v, ok := tags[l]; ok {
			ret = append(ret, v)
		} else {
			ret = append(ret, "")
		}
	}
	return ret
}