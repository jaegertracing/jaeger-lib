// Copyright (c) 2018 The Jaeger Authors
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

package metrics

// Histogram that keeps track of a distribution of values.
type Histogram interface {
	// Records the value passed in.
	Record(float64)
}

// HistogramWithExemplar that keeps track of a distribution of values and
// supports exemplars.
type HistogramWithExemplar interface {
	Histogram
	RecordWithExemplar(float64, string)
}

// NullHistogram that does nothing
var NullHistogram HistogramWithExemplar = nullHistogram{}

type nullHistogram struct{}

func (nullHistogram) Record(float64) {}

func (nullHistogram) RecordWithExemplar(float64, string) {}
