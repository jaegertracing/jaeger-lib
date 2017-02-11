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

// AssertCounterMetrics checks if counter metrics exist.
func AssertCounterMetrics(t *testing.T, f *metrics.LocalFactory, expectedMetrics ...ExpectedMetric) {
	counters, _ := f.Snapshot()
	assertMetrics(t, counters, expectedMetrics...)
}

// AssertGaugeMetrics checks if gauge metrics exist.
func AssertGaugeMetrics(t *testing.T, f *metrics.LocalFactory, expectedMetrics ...ExpectedMetric) {
	_, gauges := f.Snapshot()
	assertMetrics(t, gauges, expectedMetrics...)
}

func assertMetrics(t *testing.T, actualMetrics map[string]int64, expectedMetrics ...ExpectedMetric) {
	for _, expected := range expectedMetrics {
		key := metrics.GetKey(expected.Name, expected.Tags, "|", "=")
		assert.EqualValues(t,
			expected.Value,
			actualMetrics[key],
			"expected metric name: %s, tags: %+v", expected.Name, expected.Tags,
		)
	}
}
