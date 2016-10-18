package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitMetrics(t *testing.T) {
	testMetrics := struct {
		Gauge   Gauge   `metric:"gauge" tags:"1=one,2=two"`
		Counter Counter `metric:"counter"`
		Timer   Timer   `metric:"timer"`
	}{}

	b := NewLocalBackend(time.Minute)
	defer b.Stop()
	f := NewLocalFactory(b)

	globalTags := map[string]string{"key": "value"}

	err := initMetrics(&testMetrics, f, globalTags)
	assert.NoError(t, err)

	testMetrics.Gauge.Update(10)
	testMetrics.Counter.Inc(5)
	testMetrics.Timer.Record(time.Duration(time.Second * 35))

	// wait for metrics
	for i := 0; i < 1000; i++ {
		c, _ := b.Snapshot()
		if _, ok := c["counter"]; ok {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

	c, g := b.Snapshot()

	assert.EqualValues(t, 5, c["counter"])
	assert.EqualValues(t, 10, g["gauge"])
	assert.EqualValues(t, 36863, g["timer.P50"])

	stopwatch := StartStopwatch(testMetrics.Timer)
	stopwatch.Stop()
	assert.True(t, 0 < stopwatch.ElapsedTime())
}

var (
	noMetricTag = struct {
		NoMetricTag Counter
	}{}

	badTags = struct {
		BadTags Counter `metric:"counter" tags:"1=one,noValue"`
	}{}

	invalidMetricType = struct {
		InvalidMetricType int64 `metric:"counter"`
	}{}
)

func TestInitMetricsFailures(t *testing.T) {
	assert.EqualError(t, initMetrics(&noMetricTag, nil, nil), "Field NoMetricTag is missing a tag 'metric'")

	assert.EqualError(t, initMetrics(&badTags, nil, nil),
		"Field [BadTags]: Tag [noValue] is not of the form key=value in 'tags' string [1=one,noValue]")

	assert.EqualError(t, initMetrics(&invalidMetricType, nil, nil),
		"Field InvalidMetricType is not a pointer to timer, gauge, or counter")
}

func TestInitPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()

	Init(&noMetricTag, NullFactory, nil)
}

func TestNullMetrics(t *testing.T) {
	// This test is just for cover
	NullFactory.CreateTimer("name", nil).Record(0)
	NullFactory.CreateCounter("name", nil).Inc(0)
	NullFactory.CreateGauge("name", nil).Update(0)
}
