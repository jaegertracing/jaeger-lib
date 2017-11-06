// Copyright (c) 2017 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	promModel "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uber/jaeger-lib/metrics"
)

var _ metrics.Factory = new(Factory)

func TestCounter(t *testing.T) {
	registry := prometheus.NewPedanticRegistry()
	f1 := New(registry)
	f2 := f1.Namespace("bender", map[string]string{"a": "b"})
	c1 := f2.Counter("rodriguez", map[string]string{"x": "y"})
	c2 := f2.Counter("rodriguez", map[string]string{"x": "z"})
	c1.Inc(1)
	c1.Inc(2)
	c2.Inc(3)
	c2.Inc(4)

	snapshot, err := registry.Gather()
	require.NoError(t, err)

	m1 := findMetric(t, snapshot, "bender:rodriguez", map[string]string{"a": "b", "x": "y"})
	assert.EqualValues(t, 3, m1.GetCounter().GetValue())

	m2 := findMetric(t, snapshot, "bender:rodriguez", map[string]string{"a": "b", "x": "z"})
	assert.EqualValues(t, 7, m2.GetCounter().GetValue(), "%+v", m2)
}

func findMetric(t *testing.T, snapshot []*promModel.MetricFamily, name string, tags map[string]string) *promModel.Metric {
	for _, mf := range snapshot {
		if mf.GetName() == name {
			for _, m := range mf.GetMetric() {
				if len(m.GetLabel()) != len(tags) {
					t.Fatalf("Mismatching labels for metric %v: want %v, have %v", name, tags, m.GetLabel())
				}
				match := true
				for _, l := range m.GetLabel() {
					if v, ok := tags[l.GetName()]; !ok || v != l.GetValue() {
						match = false
					}
				}
				if match {
					return m
				}
			}
		}
	}
	t.Logf("Cannot find metric %v %v", name, tags)
	for _, nf := range snapshot {
		t.Logf("Family: %v", nf.GetName())
		for _, m := range nf.GetMetric() {
			t.Logf("==> %v", m)
		}
	}
	t.FailNow()
	return nil
}
