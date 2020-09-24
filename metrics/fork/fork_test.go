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
