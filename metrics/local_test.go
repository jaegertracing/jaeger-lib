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

	b := NewLocalBackend(time.Millisecond)
	defer b.Stop()

	f := NewLocalFactory(b)
	f.CreateCounter("my-counter", nil).Inc(4)
	f.CreateCounter("my-counter", nil).Inc(6)
	f.CreateCounter("other-counter", nil).Inc(8)
	f.CreateGauge("my-gauge", nil).Update(25)
	f.CreateGauge("my-gauge", nil).Update(43)
	f.CreateGauge("other-gauge", nil).Update(74)

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
			f.CreateTimer(metric, nil).Record(d)
		}
	}

	c, g := b.Snapshot()
	require.NotNil(t, c)
	require.NotNil(t, g)

	assert.Equal(t, map[string]int64{
		"my-counter":    10,
		"other-counter": 8,
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
