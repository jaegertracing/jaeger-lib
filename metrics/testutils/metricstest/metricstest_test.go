package metricstest

import (
	"testing"
)

func TestAssertMetrics(t *testing.T) {
	f := NewFactory(0)
	tags := map[string]string{"key": "value"}
	f.IncCounter("counter", tags, 1)
	f.UpdateGauge("gauge", tags, 11)

	AssertCounterMetrics(t, f, ExpectedMetric{Name: "counter", Tags: tags, Value: 1})
	AssertGaugeMetrics(t, f, ExpectedMetric{Name: "gauge", Tags: tags, Value: 11})
}
