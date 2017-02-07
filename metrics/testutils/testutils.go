package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uber/jaeger-lib/metrics"
)

// ExpectedMetric contains metrics under test.
type ExpectedMetric struct {
	Name  string
	Tags  map[string]string
	Value int
}

// TODO do something similar for Timers

// AssertCounterMetric checks if counter metric exists.
func AssertCounterMetric(t *testing.T, f *metrics.LocalFactory, expectedMetrics ...ExpectedMetric) {
	counters, _ := f.Snapshot()
	assertMetrics(t, counters, expectedMetrics...)
}

// AssertGaugeMetric checks if gauge metric exists.
func AssertGaugeMetric(t *testing.T, f *metrics.LocalFactory, expectedMetrics ...ExpectedMetric) {
	_, gauges := f.Snapshot()
	assertMetrics(t, gauges, expectedMetrics...)
}

// AssertCounterMetrics checks if the existing counter metrics exactly match the expected metrics.
func AssertCounterMetrics(t *testing.T, f *metrics.LocalFactory, expectedMetrics []ExpectedMetric) {
	counters, _ := f.Snapshot()
	assertMetrics(t, counters, expectedMetrics...)
	assert.Len(t, expectedMetrics, len(counters))
}

// AssertGaugeMetrics checks if the existing gauge metrics exactly match the expected metrics.
func AssertGaugeMetrics(t *testing.T, f *metrics.LocalFactory, expectedMetrics []ExpectedMetric) {
	_, gauges := f.Snapshot()
	assertMetrics(t, gauges, expectedMetrics...)
	assert.Len(t, expectedMetrics, len(gauges))
}

func assertMetrics(t *testing.T, actualMetrics map[string]int64, expectedMetrics ...ExpectedMetric) {
	for _, expected := range expectedMetrics {
		key := metrics.GetKey(expected.Name, expected.Tags)
		assert.EqualValues(t,
			expected.Value,
			actualMetrics[key],
			"expected metric name: %s, tags: %+v", expected.Name, expected.Tags,
		)
	}
}
