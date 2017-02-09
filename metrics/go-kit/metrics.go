package xkit

import (
	"time"

	kit "github.com/go-kit/kit/metrics"
)

// Counter is an adapter from go-kit Counter to jaeger-lib Counter
type Counter struct {
	counter kit.Counter
}

// NewCounter creates a new Counter
func NewCounter(counter kit.Counter) *Counter {
	return &Counter{counter: counter}
}

// Inc adds the given value to the counter.
func (c *Counter) Inc(delta int64) {
	c.counter.Add(float64(delta))
}

// Gauge is an adapter from go-kit Gauge to jaeger-lib Gauge
type Gauge struct {
	gauge kit.Gauge
}

// NewGauge creates a new Gauge
func NewGauge(gauge kit.Gauge) *Gauge {
	return &Gauge{gauge: gauge}
}

// Update the gauge to the value passed in.
func (g *Gauge) Update(value int64) {
	g.gauge.Set(float64(value))
}

// Timer is an adapter from go-kit Histogram to jaeger-lib Timer
type Timer struct {
	hist kit.Histogram
}

// NewTimer creates a new Timer
func NewTimer(hist kit.Histogram) *Timer {
	return &Timer{hist: hist}
}

// Record saves the time passed in.
func (t *Timer) Record(delta time.Duration) {
	t.hist.Observe(delta.Seconds())
}
