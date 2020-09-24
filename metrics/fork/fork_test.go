// Copyright (c) 2020 The Jaeger Authors.
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

package fork

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-lib/metrics/metricstest"
)

func TestAPICheck(t *testing.T) {
	var v interface{} = &Factory{}
	_, ok := v.(metrics.Factory)
	assert.True(t, ok)
}

func TestForkFactory(t *testing.T) {
	forkNamespace := "internal"
	forkFactory := metricstest.NewFactory(time.Second)
	defaultFactory := metricstest.NewFactory(time.Second)

	// Create factory that will delegate namespaced metrics to forkFactory
	// and add some metrics
	ff := New(forkNamespace, forkFactory, defaultFactory)
	ff.Gauge(metrics.Options{
		Name: "somegauge",
		Tags: nil,
		Help: "",
	}).Update(666)
	ff.Counter(metrics.Options{
		Name: "somecounter",
		Tags: nil,
		Help: "",
	}).Inc(2)

	// Check that metrics are presented in defaultFactory backend
	defaultFactory.AssertCounterMetrics(t, metricstest.ExpectedMetric{
		Name:  "somecounter",
		Tags:  nil,
		Value: 2,
	})
	defaultFactory.AssertGaugeMetrics(t, metricstest.ExpectedMetric{
		Name:  "somegauge",
		Tags:  nil,
		Value: 666,
	})

	// Get default namespaced factory
	defaultNamespacedFactory := ff.Namespace(metrics.NSOptions{
		Name: "default",
		Tags: nil,
	})

	// Add some metrics
	defaultNamespacedFactory.Counter(metrics.Options{
		Name: "somenamespacedcounter",
		Tags: nil,
		Help: "",
	}).Inc(111)
	defaultNamespacedFactory.Gauge(metrics.Options{
		Name: "somenamespacedgauge",
		Tags: nil,
		Help: "",
	}).Update(222)
	defaultNamespacedFactory.Histogram(metrics.HistogramOptions{
		Name:    "somenamespacedhist",
		Tags:    nil,
		Help:    "",
		Buckets: nil,
	}).Record(1)
	defaultNamespacedFactory.Timer(metrics.TimerOptions{
		Name:    "somenamespacedtimer",
		Tags:    nil,
		Help:    "",
		Buckets: nil,
	}).Record(time.Millisecond)

	// Check values in default namespaced factory backend
	defaultFactory.AssertCounterMetrics(t, metricstest.ExpectedMetric{
		Name:  "default.somenamespacedcounter",
		Tags:  nil,
		Value: 111,
	})
	defaultFactory.AssertGaugeMetrics(t, metricstest.ExpectedMetric{
		Name:  "default.somenamespacedgauge",
		Tags:  nil,
		Value: 222,
	})

	// Get factory with forkNamespace and add some metrics
	internalFactory := ff.Namespace(metrics.NSOptions{
		Name: forkNamespace,
		Tags: nil,
	})
	internalFactory.Gauge(metrics.Options{
		Name: "someinternalgauge",
		Tags: nil,
		Help: "",
	}).Update(20)
	internalFactory.Counter(metrics.Options{
		Name: "someinternalcounter",
		Tags: nil,
		Help: "",
	}).Inc(50)

	// Check that metrics are presented in forkFactory backend
	forkFactory.AssertGaugeMetrics(t, metricstest.ExpectedMetric{
		Name:  "internal.someinternalgauge",
		Tags:  nil,
		Value: 20,
	})
	forkFactory.AssertCounterMetrics(t, metricstest.ExpectedMetric{
		Name:  "internal.someinternalcounter",
		Tags:  nil,
		Value: 50,
	})
}
