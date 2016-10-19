package metrics

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/codahale/hdrhistogram"
)

const (
	defaultCollectionInterval = time.Minute
)

// TODO tags are currently ignored by LocalBackend, this should be fixed for local
// metric tests that are sensitive to the tags

// This is intentionally very similar to github.com/codahale/metrics, the
// main difference being that counters/gauges are scoped to the provider
// rather than being global (to facilitate testing).

// A LocalBackend is a metrics provider which aggregates data in-vm, and
// allows exporting snapshots to shove the data into a remote collector
type LocalBackend struct {
	cm       sync.RWMutex
	gm       sync.RWMutex
	tm       sync.RWMutex
	counters map[string]*int64
	gauges   map[string]*int64
	timers   map[string]*localBackendTimer
	stop     chan struct{}
	stopped  chan struct{}
}

// NewLocalBackend returns a new LocalBackend. The collectionInterval is the histogram
// time window for each timer.
func NewLocalBackend(collectionInterval time.Duration) *LocalBackend {
	b := &LocalBackend{
		counters: make(map[string]*int64),
		gauges:   make(map[string]*int64),
		timers:   make(map[string]*localBackendTimer),
		stop:     make(chan struct{}),
		stopped:  make(chan struct{}),
	}

	if collectionInterval == 0 {
		collectionInterval = defaultCollectionInterval
	}

	go b.runLoop(collectionInterval)
	return b
}

func (b *LocalBackend) runLoop(collectionInterval time.Duration) {
	ticker := time.NewTicker(collectionInterval)
	for {
		select {
		case <-ticker.C:
			b.tm.Lock()
			timers := make(map[string]*localBackendTimer, len(b.timers))
			for timerName, timer := range b.timers {
				timers[timerName] = timer
			}
			b.tm.Unlock()

			for _, t := range timers {
				t.Lock()
				t.hist.Rotate()
				t.Unlock()
			}
		case <-b.stop:
			ticker.Stop()
			close(b.stopped)
			return
		}
	}
}

// IncCounter increments a counter value
func (b *LocalBackend) IncCounter(name string, tags map[string]string, delta int64) {
	name = b.getKey(name, tags)
	b.cm.RLock()
	counter := b.counters[name]
	b.cm.RUnlock()

	if counter != nil {
		atomic.AddInt64(counter, delta)
		return
	}

	b.cm.Lock()
	counter = b.counters[name]
	if counter == nil {
		b.counters[name] = new(int64)
		*b.counters[name] = delta
	} else {
		atomic.AddInt64(counter, delta)
	}
	b.cm.Unlock()
}

// UpdateGauge updates the value of a gauge
func (b *LocalBackend) UpdateGauge(name string, tags map[string]string, value int64) {
	name = b.getKey(name, tags)
	b.gm.RLock()
	gauge := b.gauges[name]
	b.gm.RUnlock()

	if gauge != nil {
		atomic.StoreInt64(gauge, value)
		return
	}

	b.gm.Lock()
	gauge = b.gauges[name]
	if gauge == nil {
		b.gauges[name] = new(int64)
		*b.gauges[name] = value
	} else {
		atomic.StoreInt64(gauge, value)
	}
	b.gm.Unlock()
}

// RecordTimer records a timing duration
func (b *LocalBackend) RecordTimer(name string, tags map[string]string, d time.Duration) {
	name = b.getKey(name, tags)
	timer := b.findOrCreateTimer(name)
	timer.Lock()
	timer.hist.Current.RecordValue(int64(d / time.Millisecond))
	timer.Unlock()
}

func (b *LocalBackend) findOrCreateTimer(name string) *localBackendTimer {
	b.tm.RLock()
	t := b.timers[name]
	b.tm.RUnlock()

	if t != nil {
		return t
	}

	b.tm.Lock()
	defer b.tm.Unlock()
	t = b.timers[name]
	if t != nil {
		return t
	}

	t = &localBackendTimer{
		hist: hdrhistogram.NewWindowed(5, 0, int64((5*time.Minute)/time.Millisecond), 1),
	}
	b.timers[name] = t
	return t
}

type localBackendTimer struct {
	sync.Mutex
	hist *hdrhistogram.WindowedHistogram
}

var (
	percentiles = map[string]float64{
		"P50":  50,
		"P75":  75,
		"P90":  90,
		"P95":  95,
		"P99":  99,
		"P999": 99.9,
	}
)

// Snapshot captures a snapshot of the current counter and gauge values
func (b *LocalBackend) Snapshot() (counters, gauges map[string]int64) {
	b.cm.RLock()
	defer b.cm.RUnlock()

	counters = make(map[string]int64, len(b.counters))
	for name, value := range b.counters {
		counters[name] = atomic.LoadInt64(value)
	}

	b.gm.RLock()
	defer b.gm.RUnlock()

	gauges = make(map[string]int64, len(b.gauges))
	for name, value := range b.gauges {
		gauges[name] = atomic.LoadInt64(value)
	}

	b.tm.RLock()
	timers := make(map[string]*localBackendTimer)
	for timerName, timer := range b.timers {
		timers[timerName] = timer
	}
	b.tm.RUnlock()

	for timerName, timer := range timers {
		timer.Lock()
		hist := timer.hist.Merge()
		timer.Unlock()
		for name, q := range percentiles {
			gauges[timerName+"."+name] = hist.ValueAtQuantile(q)
		}
	}

	return
}

// Stop cleanly closes the background goroutine spawned by NewLocalBackend.
func (b *LocalBackend) Stop() {
	close(b.stop)
	<-b.stopped
}

// MetricDescr describes a metric with tags
type MetricDescr struct {
	Name string
	Tags map[string]string
}

// Key converts name+tags into a single string of the form
// "name|tag1=value1|...|tagN=valueN", where tag names are
// sorted alphabetically.
func (m MetricDescr) Key() string {
	var keys []string
	for k := range m.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	key := m.Name
	for _, k := range keys {
		key = key + "|" + k + "=" + m.Tags[k]
	}
	return key
}

func (b *LocalBackend) getKey(name string, tags map[string]string) string {
	return MetricDescr{
		Name: name,
		Tags: tags,
	}.Key()
}

type stats struct {
	name         string
	tags         map[string]string
	localBackend *LocalBackend
}

type localTimer struct {
	stats
}

func (l *localTimer) Record(d time.Duration) {
	l.localBackend.RecordTimer(l.name, l.tags, d)
}

type localCounter struct {
	stats
}

func (l *localCounter) Inc(delta int64) {
	l.localBackend.IncCounter(l.name, l.tags, delta)
}

type localGauge struct {
	stats
}

func (l *localGauge) Update(value int64) {
	l.localBackend.UpdateGauge(l.name, l.tags, value)
}

// LocalFactory stats factory that creates metrics that are stored locally
type LocalFactory struct {
	localBackend *LocalBackend
}

// NewLocalFactory returns a new LocalMetricsFactory
func NewLocalFactory(lb *LocalBackend) Factory {
	return &LocalFactory{
		localBackend: lb,
	}
}

// CreateCounter returns a local stats counter
func (l *LocalFactory) CreateCounter(name string, tags map[string]string) Counter {
	return &localCounter{
		stats{
			name:         name,
			tags:         tags,
			localBackend: l.localBackend,
		},
	}
}

// CreateTimer returns a local stats timer
func (l *LocalFactory) CreateTimer(name string, tags map[string]string) Timer {
	return &localTimer{
		stats{
			name:         name,
			tags:         tags,
			localBackend: l.localBackend,
		},
	}
}

// CreateGauge returns a local stats gauge
func (l *LocalFactory) CreateGauge(name string, tags map[string]string) Gauge {
	return &localGauge{
		stats{
			name:         name,
			tags:         tags,
			localBackend: l.localBackend,
		},
	}
}
