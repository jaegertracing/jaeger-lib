package metrics

import (
	"testing"
)

func TestAssertMetrics(t *testing.T) {
	b := NewLocalBackend(0)
	tags := map[string]string{"key": "value"}
	b.IncCounter("counter", tags, 1)
	b.UpdateGauge("gauge", tags, 11)

	AssertCounterMetrics(t, b, []ExpectedMetric{{Name: "counter", Tags: tags, Value: 1}})
	AssertGaugeMetrics(t, b, []ExpectedMetric{{Name: "gauge", Tags: tags, Value: 11}})
}
