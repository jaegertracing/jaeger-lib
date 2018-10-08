package testutils

import (
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"github.com/uber/jaeger-lib/metrics/expvar"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	tally2 "github.com/uber/jaeger-lib/metrics/tally"
	"testing"

	"github.com/uber/jaeger-lib/metrics"
)

func TestAssertMetrics(t *testing.T) {
	f := metrics.NewLocalFactory(0)
	tags := map[string]string{"key": "value"}
	f.IncCounter("counter", tags, 1)
	f.UpdateGauge("gauge", tags, 11)

	AssertCounterMetrics(t, f, ExpectedMetric{Name: "counter", Tags: tags, Value: 1})
	AssertGaugeMetrics(t, f, ExpectedMetric{Name: "gauge", Tags: tags, Value: 11})
}

func TestNameCaching(t *testing.T) {
	tags := map[string]string{
		"partition": "12",
	}
	metricName := "start-consuming"

	// Expvar
	ef := expvar.NewFactory(2)
	ec1 := ef.Counter(metricName, tags)
	ec2 := ef.Counter(metricName, tags)
	assert.Equal(t, ec1, ec2)

	// Prom
	pf := prometheus.New()
	pc1 := pf.Counter(metricName, tags)
	pc2 := pf.Counter(metricName, tags)
	assert.Equal(t, pc1, pc2)

	// Tally
	ts, _ := tally.NewRootScope(tally.ScopeOptions{}, 0)
	tf := tally2.Wrap(ts)
	tc1 := tf.Counter(metricName, tags)
	tc2 := tf.Counter(metricName, tags)
	assert.Equal(t, tc1, tc2)
}
