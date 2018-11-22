// Copyright (c) 2017 Uber Technologies, Inc.
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

package tally

import (
	"github.com/uber-go/tally"

	"github.com/uber/jaeger-lib/metrics"
)

// Wrap takes a tally Scope and returns jaeger-lib metrics.Factory.
func Wrap(scope tally.Scope) metrics.Factory {
	return &factory{
		tally: scope,
	}
}

// TODO implement support for tags if tally.Scope does not support them
type factory struct {
	tally tally.Scope
}

func (f *factory) Counter(metricInfo metrics.MetricInfo) metrics.Counter {
	scope := f.tally
	if len(metricInfo.Tags) > 0 {
		scope = scope.Tagged(metricInfo.Tags)
	}
	return NewCounter(scope.Counter(metricInfo.Name))
}

func (f *factory) Gauge(metricInfo metrics.MetricInfo) metrics.Gauge {
	scope := f.tally
	if len(metricInfo.Tags) > 0 {
		scope = scope.Tagged(metricInfo.Tags)
	}
	return NewGauge(scope.Gauge(metricInfo.Name))
}

func (f *factory) Timer(metricInfo metrics.MetricInfo) metrics.Timer {
	scope := f.tally
	if len(metricInfo.Tags) > 0 {
		scope = scope.Tagged(metricInfo.Tags)
	}
	return NewTimer(scope.Timer(metricInfo.Name))
}

func (f *factory) Namespace(metricScope metrics.MetricScope) metrics.Factory {
	return &factory{
		tally: f.tally.SubScope(metricScope.Name).Tagged(metricScope.Tags),
	}
}
