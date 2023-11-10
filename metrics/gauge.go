package metrics

import "sync/atomic"

// gaugeSnapshot contains a readonly int64.
type GaugeSnapshot interface {
	Value() int64
}

// Gauges hold an int64 value that can be set arbitrarily.
type Gauge interface {
	Snapshot() GaugeSnapshot
	Update(int64)
	UpdateIfGt(int64)
	Dec(int64)
	Inc(int64)
}

// GetOrRegisterGauge returns an existing Gauge or constructs and registers a
// new StandardGauge.
func GetOrRegisterGauge(name string, r Registry) Gauge {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewGauge).(Gauge)
}

// NewGauge constructs a new StandardGauge.
func NewGauge() Gauge {
	if !Enabled {
		return NilGauge{}
	}
	return &StandardGauge{}
}

// NewRegisteredGauge constructs and registers a new StandardGauge.
func NewRegisteredGauge(name string, r Registry) Gauge {
	c := NewGauge()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NewFunctionalGauge constructs a new FunctionalGauge.
func NewFunctionalGauge(f func() int64) Gauge {
	if !Enabled {
		return NilGauge{}
	}
	return &FunctionalGauge{value: f}
}

// NewRegisteredFunctionalGauge constructs and registers a new StandardGauge.
func NewRegisteredFunctionalGauge(name string, r Registry, f func() int64) Gauge {
	c := NewFunctionalGauge(f)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// gaugeSnapshot is a read-only copy of another Gauge.
type gaugeSnapshot int64

// Value returns the value at the time the snapshot was taken.
func (g gaugeSnapshot) Value() int64 { return int64(g) }

// NilGauge is a no-op Gauge.
type NilGauge struct{}

func (NilGauge) Snapshot() GaugeSnapshot { return (*emptySnapshot)(nil) }
func (NilGauge) Update(v int64)          {}
func (NilGauge) UpdateIfGt(v int64)      {}
func (NilGauge) Dec(i int64)             {}
func (NilGauge) Inc(i int64)             {}

// StandardGauge is the standard implementation of a Gauge and uses the
// sync/atomic package to manage a single int64 value.
type StandardGauge struct {
	value atomic.Int64
}

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGauge) Snapshot() GaugeSnapshot {
	return gaugeSnapshot(g.value.Load())
}

// Update updates the gauge's value.
func (g *StandardGauge) Update(v int64) {
	g.value.Store(v)
}

// Update updates the gauge's value if v is larger then the current valie.
func (g *StandardGauge) UpdateIfGt(v int64) {
	for {
		exist := g.value.Load()
		if exist >= v {
			break
		}
		if g.value.CompareAndSwap(exist, v) {
			break
		}
	}
}

// Dec decrements the gauge's current value by the given amount.
func (g *StandardGauge) Dec(i int64) {
	g.value.Add(-i)
}

// Inc increments the gauge's current value by the given amount.
func (g *StandardGauge) Inc(i int64) {
	g.value.Add(i)
}

// FunctionalGauge returns value from given function
type FunctionalGauge struct {
	value func() int64
}

func (g FunctionalGauge) UpdateIfGt(i int64) {
	//TODO implement me
	panic("implement me")
}

// Value returns the gauge's current value.
func (g FunctionalGauge) Value() int64 {
	return g.value()
}

// Snapshot returns the snapshot.
func (g FunctionalGauge) Snapshot() GaugeSnapshot { return gaugeSnapshot(g.Value()) }

// Update panics.
func (FunctionalGauge) Update(int64) {
	panic("Update called on a FunctionalGauge")
}

// Dec panics.
func (FunctionalGauge) Dec(int64) {
	panic("Dec called on a FunctionalGauge")
}

// Inc panics.
func (FunctionalGauge) Inc(int64) {
	panic("Inc called on a FunctionalGauge")
}
