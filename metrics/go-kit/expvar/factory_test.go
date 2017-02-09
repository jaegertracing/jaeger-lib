package expvar

import (
	"expvar"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	f := NewFactory("gokit_expvar", 10)
	c := f.Counter("counter", nil)
	c.Inc(42)
	kv := findExpvar("gokit_expvar_counter")
	assert.Equal(t, "42", kv.Value.String())
}

func TestGauge(t *testing.T) {
	f := NewFactory("gokit_expvar", 10)
	g := f.Gauge("gauge", nil)
	g.Update(42)
	kv := findExpvar("gokit_expvar_gauge")
	assert.Equal(t, "42", kv.Value.String())
}

func TestTimer(t *testing.T) {
	f := NewFactory("gokit_expvar", 10)
	timer := f.Timer("timer", nil)
	timer.Record(100*time.Millisecond + 500*time.Microsecond)
	kv := findExpvar("gokit_expvar_timer.p50")
	assert.Equal(t, "100.5", kv.Value.String())
}

func TestNamespace(t *testing.T) {
	f := NewFactory("gokit_expvar", 10)
	f = f.Namespace("namespace", nil)
	c := f.Counter("counter", nil)
	c.Inc(42)
	kv := findExpvar("gokit_expvar_namespace_counter")
	assert.Equal(t, "42", kv.Value.String())
}

func findExpvar(key string) *expvar.KeyValue {
	var kv *expvar.KeyValue
	expvar.Do(func(v expvar.KeyValue) {
		if v.Key == key {
			kv = &v
		}
	})
	return kv
}
