package tally

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/uber-go/tally"
)

func TestFactory(t *testing.T) {
	testScope := tally.NewTestScope("pre", map[string]string{"a": "b"})
	factory := Wrap(testScope).Namespace("fix", map[string]string{"c": "d"})
	counter := factory.Counter("counter", map[string]string{"x": "y"})
	counter.Inc(42)
	gauge := factory.Gauge("gauge", map[string]string{"x": "y"})
	gauge.Update(42)
	timer := factory.Timer("timer", map[string]string{"x": "y"})
	timer.Record(42 * time.Millisecond)
	snapshot := testScope.Snapshot()
	c := snapshot.Counters()["pre.fix.counter"]
	g := snapshot.Gauges()["pre.fix.gauge"]
	h := snapshot.Timers()["pre.fix.timer"]
	expectedTags := map[string]string{"a": "b", "c": "d", "x": "y"}
	assert.EqualValues(t, 42, c.Value())
	assert.EqualValues(t, expectedTags, c.Tags())
	assert.EqualValues(t, 42, g.Value())
	assert.EqualValues(t, expectedTags, g.Tags())
	assert.Equal(t, []time.Duration{42 * time.Millisecond}, h.Values())
	assert.EqualValues(t, expectedTags, h.Tags())
}
