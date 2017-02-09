package xkit

import (
	"testing"
	"time"

	kit "github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/generic"
	"github.com/stretchr/testify/assert"

	"github.com/uber/jaeger-lib/metrics"
)

type genericFactory struct{}

func (f genericFactory) Counter(name string) kit.Counter     { return generic.NewCounter(name) }
func (f genericFactory) Gauge(name string) kit.Gauge         { return generic.NewGauge(name) }
func (f genericFactory) Histogram(name string) kit.Histogram { return generic.NewHistogram(name, 10) }

type Tags map[string]string
type metricFunc func(t *testing.T, testCase testCase, f metrics.Factory) (name func() string, labels func() []string)

type testCase struct {
	prefix string
	name   string
	tags   Tags

	useNamespace  bool
	namespace     string
	namespaceTags Tags

	expName string
	expTags []string
}

func TestFactoryScoping(t *testing.T) {
	testSuites := []struct {
		metricType string
		metricFunc metricFunc
	}{
		{"counter", testCounter},
		{"gauge", testGauge},
		{"timer", testTimer},
	}
	for _, ts := range testSuites {
		testSuite := ts // capture loop var
		testCases := []testCase{
			{prefix: "x", name: "", expName: "x"},
			{prefix: "", name: "y", expName: "y"},
			{prefix: "x", name: "y", expName: "x.y"},
			{prefix: "x", name: "z", expName: "x.z", tags: Tags{"a": "b"}, expTags: []string{"a", "b"}},
			{
				name:         "y",
				useNamespace: true,
				namespace:    "z",
				expName:      "z.y",
			},
			{
				name:          "y",
				useNamespace:  true,
				namespace:     "w",
				namespaceTags: Tags{"a": "b"},
				expName:       "w.y",
				expTags:       []string{"a", "b"},
			},
		}
		for _, tc := range testCases {
			testCase := tc // capture loop var
			t.Run(testSuite.metricType+":"+testCase.expName, func(t *testing.T) {
				f := Wrap(testCase.prefix, genericFactory{})
				if testCase.useNamespace {
					f = f.Namespace(testCase.namespace, testCase.namespaceTags)
				}
				name, labels := testSuite.metricFunc(t, testCase, f)
				if testCase.tags == nil && testCase.namespaceTags == nil {
					// TODO go-kit loses the name of the counter on With()
					// https://github.com/go-kit/kit/issues/455
					assert.Equal(t, testCase.expName, name())
				}
				if testCase.expTags != nil {
					assert.Equal(t, testCase.expTags, labels())
				}
			})
		}
	}
}

func testCounter(t *testing.T, testCase testCase, f metrics.Factory) (name func() string, labels func() []string) {
	c := f.Counter(testCase.name, testCase.tags)
	c.Inc(123)
	gc := c.(*Counter).counter.(*generic.Counter)
	assert.EqualValues(t, 123.0, gc.Value())
	name = func() string { return gc.Name }
	labels = gc.LabelValues
	return
}

func testGauge(t *testing.T, testCase testCase, f metrics.Factory) (name func() string, labels func() []string) {
	g := f.Gauge(testCase.name, testCase.tags)
	g.Update(123)
	gg := g.(*Gauge).gauge.(*generic.Gauge)
	assert.EqualValues(t, 123.0, gg.Value())
	name = func() string { return gg.Name }
	labels = gg.LabelValues
	return
}

func testTimer(t *testing.T, testCase testCase, f metrics.Factory) (name func() string, labels func() []string) {
	tm := f.Timer(testCase.name, testCase.tags)
	tm.Record(123 * time.Millisecond)
	gt := tm.(*Timer).hist.(*generic.Histogram)
	assert.InDelta(t, 0.123, gt.Quantile(0.9), 0.00001)
	name = func() string { return gt.Name }
	labels = gt.LabelValues
	return
}
