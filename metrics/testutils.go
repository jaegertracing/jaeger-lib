package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ExpectedMetric contains metrics under test.
type ExpectedMetric struct {
	Name  string
	Tags  map[string]string
	Value int
}

// TODO do something similar for Timers

// AssertCounterMetrics checks if counter metrics exist.
func AssertCounterMetrics(t *testing.T, b *LocalBackend, expectedMetrics []ExpectedMetric) {
	counters, _ := b.Snapshot()
	assertMetrics(t, counters, expectedMetrics)
}

// AssertGaugeMetrics checks if gauge metrics exist.
func AssertGaugeMetrics(t *testing.T, b *LocalBackend, expectedMetrics []ExpectedMetric) {
	_, gauges := b.Snapshot()
	assertMetrics(t, gauges, expectedMetrics)
}

func assertMetrics(t *testing.T, actualMetrics map[string]int64, expectedMetrics []ExpectedMetric) {
	for _, expected := range expectedMetrics {
		key := getKey(expected.Name, expected.Tags)
		assert.EqualValues(t,
			expected.Value,
			actualMetrics[key],
			"expected metric name: %s, tags: %+v", expected.Name, expected.Tags,
		)
	}
	assert.Len(t, expectedMetrics, len(actualMetrics))
}
