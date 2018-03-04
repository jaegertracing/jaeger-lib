// Copyright (c) 2018 Uber Technologies, Inc.
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

package adapters

import (
	"sync"

	"github.com/uber/jaeger-lib/metrics"
)

type cache struct {
	lock     sync.Mutex
	counters map[string]metrics.Counter
	gauges   map[string]metrics.Gauge
	timers   map[string]metrics.Timer
}

func newCache() *cache {
	return &cache{
		counters: make(map[string]metrics.Counter),
		gauges:   make(map[string]metrics.Gauge),
		timers:   make(map[string]metrics.Timer),
	}
}

func (r *cache) getCounter(name string) (metrics.Counter, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	c, ok := r.counters[name]
	return c, ok
}

func (r *cache) getGauge(name string) (metrics.Gauge, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	g, ok := r.gauges[name]
	return g, ok
}

func (r *cache) getTimer(name string) (metrics.Timer, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	t, ok := r.timers[name]
	return t, ok
}

func (r *cache) setCounter(name string, c metrics.Counter) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.counters[name] = c
}

func (r *cache) setGauge(name string, g metrics.Gauge) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.gauges[name] = g
}

func (r *cache) setTimer(name string, t metrics.Timer) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.timers[name] = t
}
