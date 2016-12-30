package metrics

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalMetrics(t *testing.T) {
	numGoroutines := runtime.NumGoroutine()
	defer func() {
		assert.Equal(t, numGoroutines, runtime.NumGoroutine(), "Leaked at least one goroutine.")
	}()

	b := NewLocalBackend(0) // default interval
	defer b.Stop()

	tags := map[string]string{
		"x": "y",
	}

	f := NewLocalFactory(b)
	f.Counter("my-counter", tags).Inc(4)
	f.Counter("my-counter", nil).Inc(6)
	f.Counter("other-counter", nil).Inc(8)
	f.Gauge("my-gauge", nil).Update(25)
	f.Gauge("my-gauge", nil).Update(43)
	f.Gauge("other-gauge", nil).Update(74)
	f.Namespace("namespace", tags).Counter("my-counter", nil).Inc(7)

	timings := map[string][]time.Duration{
		"foo-latency": {
			time.Second * 35,
			time.Second * 6,
			time.Millisecond * 576,
			time.Second * 12,
		},
		"bar-latency": {
			time.Minute*4 + time.Second*34,
			time.Minute*7 + time.Second*12,
			time.Second * 625,
			time.Second * 12,
		},
	}

	for metric, timing := range timings {
		for _, d := range timing {
			f.Timer(metric, nil).Record(d)
		}
	}

	c, g := b.Snapshot()
	require.NotNil(t, c)
	require.NotNil(t, g)

	assert.Equal(t, map[string]int64{
		"my-counter":           10,
		"other-counter":        8,
		"namespace.my-counter": 7,
	}, c)

	assert.Equal(t, map[string]int64{
		"bar-latency.P50":  278527,
		"bar-latency.P75":  278527,
		"bar-latency.P90":  442367,
		"bar-latency.P95":  442367,
		"bar-latency.P99":  442367,
		"bar-latency.P999": 442367,
		"foo-latency.P50":  6143,
		"foo-latency.P75":  12287,
		"foo-latency.P90":  36863,
		"foo-latency.P95":  36863,
		"foo-latency.P99":  36863,
		"foo-latency.P999": 36863,
		"my-gauge":         43,
		"other-gauge":      74,
	}, g)
}

func TestLocalMetricsInterval(t *testing.T) {
	numGoroutines := runtime.NumGoroutine()
	defer func() {
		assert.Equal(t, numGoroutines, runtime.NumGoroutine(), "Leaked at least one goroutine.")
	}()

	refreshInterval := time.Millisecond
	const relativeCheckFrequency = 5 // check 5 times per refreshInterval
	const maxChecks = 2 * relativeCheckFrequency
	checkInterval := (refreshInterval * relativeCheckFrequency) / maxChecks

	b := NewLocalBackend(refreshInterval)
	defer b.Stop()

	f := NewLocalFactory(b)
	f.Timer("timer", nil).Record(1)

	b.tm.Lock()
	timer := b.timers["timer"]
	b.tm.Unlock()
	assert.NotNil(t, timer)

	// timer.hist.Current is modified on every Rotate(), which is called by LocalBackend after every refreshInterval
	getCurr := func() interface{} {
		timer.Lock()
		defer timer.Unlock()
		return timer.hist.Current
	}

	curr := getCurr()

	// wait for twice as long as the refresh interval
	for i := 0; i < maxChecks; i++ {
		time.Sleep(checkInterval)

		if getCurr() != curr {
			return
		}
	}
	t.Fail()
}
