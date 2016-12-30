package metrics

// Factory creates new metrics
type Factory interface {
	Counter(name string, tags map[string]string) Counter
	Timer(name string, tags map[string]string) Timer
	Gauge(name string, tags map[string]string) Gauge

	// Namespace returns a nested metrics factory.
	Namespace(name string, tags map[string]string) Factory
}

// NullFactory is a metrics factory that returns NullCounter, NullTimer, and NullGauge.
var NullFactory Factory = nullFactory{}

type nullFactory struct{}

func (nullFactory) Counter(name string, tags map[string]string) Counter   { return NullCounter }
func (nullFactory) Timer(name string, tags map[string]string) Timer       { return NullTimer }
func (nullFactory) Gauge(name string, tags map[string]string) Gauge       { return NullGauge }
func (nullFactory) Namespace(name string, tags map[string]string) Factory { return NullFactory }
