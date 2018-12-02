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

// NSOptions defines the name and tags map associated with a metric
type NSOptions struct {
	Name string
	Tags map[string]string
}

// Options defines the information associated with a metric
type Options struct {
	Name string
	Tags map[string]string
	Help string
}

// Factory creates new metrics
type Factory interface {
	Counter(metric Options) Counter
	Timer(metric Options) Timer
	Gauge(metric Options) Gauge

	// Namespace returns a nested metrics factory.
	Namespace(scope NSOptions) Factory
}

// NullFactory is a metrics factory that returns NullCounter, NullTimer, and NullGauge.
var NullFactory Factory = nullFactory{}

type nullFactory struct{}

func (nullFactory) Counter(options Options) Counter {
	return NullCounter
}
func (nullFactory) Timer(options Options) Timer {
	return NullTimer
}
func (nullFactory) Gauge(options Options) Gauge {
	return NullGauge
}
func (nullFactory) Namespace(scope NSOptions) Factory { return NullFactory }
