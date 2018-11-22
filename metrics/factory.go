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

package metrics

// MetricScope defines the name and tags map associated with a metric
type MetricScope struct {
	Name string
	Tags map[string]string
}

// MetricInfo defines the information associated with a metric
type MetricInfo struct {
	MetricScope
	Description string
}

// Factory creates new metrics
type Factory interface {
	Counter(metric MetricInfo) Counter
	Timer(metric MetricInfo) Timer
	Gauge(metric MetricInfo) Gauge

	// Namespace returns a nested metrics factory.
	Namespace(metricScope MetricScope) Factory
}

// NullFactory is a metrics factory that returns NullCounter, NullTimer, and NullGauge.
var NullFactory Factory = nullFactory{}

type nullFactory struct{}

func (nullFactory) Counter(metricInfo MetricInfo) Counter {
	return NullCounter
}
func (nullFactory) Timer(metricInfo MetricInfo) Timer {
	return NullTimer
}
func (nullFactory) Gauge(metricInfo MetricInfo) Gauge {
	return NullGauge
}
func (nullFactory) Namespace(metricScope MetricScope) Factory { return NullFactory }
