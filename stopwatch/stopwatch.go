// Tideland Go Trace - Stopwatch
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stopwatch // import "tideland.dev/go/trace/stopwatch"

//--------------------
// IMPORT
//--------------------

import (
	"fmt"
	"sync"
	"time"
)

//--------------------
// MEASUREMENT
//--------------------

// Measurement contains one single measurement.
type Measurement struct {
	owner *MeteringPoint
	start time.Time
}

// Stop ends the measurement and enques its duration.
func (m *Measurement) Stop() {
	m.owner.enqueue(time.Since(m.start))
}

//--------------------
// METERING POINT VALUE
//--------------------

// MeteringPointValue contains the accumulated value of one metering point.
type MeteringPointValue struct {
	Namespace string        `json:"namespace"`
	ID        string        `json:"id"`
	Quantity  int           `json:"quantity"`
	Total     time.Duration `json:"total"`
	Minimum   time.Duration `json:"minimum"`
	Maximum   time.Duration `json:"maximum"`
	Average   time.Duration `json:"average"`
}

// String implements the fmt.Stringer interface.
func (mpv MeteringPointValue) String() string {
	return fmt.Sprintf(
		"[%s :: %s] %d / %v / %v / %v / %v",
		mpv.Namespace,
		mpv.ID,
		mpv.Quantity,
		mpv.Total,
		mpv.Minimum,
		mpv.Maximum,
		mpv.Average,
	)
}

// MeteringPointValues contains a set of accumulated metering point values.
type MeteringPointValues []MeteringPointValue

//--------------------
// METERING POINT
//--------------------

// MeteringPoint collects the measurements of one code section.
type MeteringPoint struct {
	mu       sync.Mutex
	owner    *Stopwatch
	id       string
	queue    []time.Duration
	quantity int
	total    time.Duration
	minimum  time.Duration
	maximum  time.Duration
}

// Start begins a new measurement.
func (mp *MeteringPoint) Start() Measurement {
	return Measurement{
		owner: mp,
		start: time.Now(),
	}
}

// Measure is a convenience function to measure one anonymous function.
func (mp *MeteringPoint) Measure(f func()) {
	m := mp.Start()
	defer m.Stop()
	f()
}

// Value returns the current value after an accumulation.
func (mp *MeteringPoint) Value() MeteringPointValue {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.accumulate()
	return MeteringPointValue{
		Namespace: mp.owner.namespace,
		ID:        mp.id,
		Quantity:  mp.quantity,
		Total:     mp.total,
		Minimum:   mp.minimum,
		Maximum:   mp.maximum,
		Average:   time.Duration(int64(mp.total) / int64(mp.quantity)),
	}
}

// enqueue adds the measured duration.
func (mp *MeteringPoint) enqueue(measuring time.Duration) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	if len(mp.queue) == cap(mp.queue) {
		// Accumulate the enqueued values.
		mp.accumulate()
	}
	mp.queue = append(mp.queue, measuring)
}

// accumulate evaluates the collected measurings synchronous or asynchronous.
func (mp *MeteringPoint) accumulate() {
	// Accumulate set of enqueued measurings.
	for _, duration := range mp.queue {
		if mp.quantity == 0 {
			mp.quantity = 1
			mp.total = duration
			mp.minimum = duration
			mp.maximum = duration
			continue
		}
		mp.quantity++
		mp.total += duration
		if mp.minimum > duration {
			mp.minimum = duration
		}
		if mp.maximum < duration {
			mp.maximum = duration
		}
	}
	mp.clearQueue()
}

// clearQueue clears the measuring queue.
func (mp *MeteringPoint) clearQueue() {
	mp.queue = make([]time.Duration, 0, 1024)
}

//--------------------
// STOPWATCH
//--------------------

// Stopwatch allows to measure the execution time at multiple reading points
// in one namespace.
type Stopwatch struct {
	mu             sync.RWMutex
	namespace      string
	meteringPoints map[string]*MeteringPoint
}

// MeteringPoint returns a new or already existing metering point
// with the given ID.
func (sw *Stopwatch) MeteringPoint(id string) *MeteringPoint {
	// First check for existing metering point.
	sw.mu.RLock()
	mp, ok := sw.meteringPoints[id]
	sw.mu.RUnlock()
	if ok {
		return mp
	}
	// Not yet existing.
	sw.mu.Lock()
	mp = &MeteringPoint{
		owner: sw,
		id:    id,
	}
	mp.clearQueue()
	sw.meteringPoints[id] = mp
	sw.mu.Unlock()
	return mp
}

// Values returns the accumulated metering point values of this stopwatch.
func (sw *Stopwatch) Values() MeteringPointValues {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	mpvs := MeteringPointValues{}
	for _, mp := range sw.meteringPoints {
		mpvs = append(mpvs, mp.Value())
	}
	return mpvs
}

//--------------------
// REGISTRY
//--------------------

// Registry is the register type for all Stopwatches.
type Registry struct {
	mu      sync.RWMutex
	watches map[string]*Stopwatch
}

// New returns a registry for Stopwatches.
func New() *Registry {
	r := &Registry{
		watches: make(map[string]*Stopwatch),
	}
	return r
}

// ForNamespace creates an instance of a Stopwatch with the given namespace.
// In case that namespace is already in use that Stopwatch will be returned.
func (r *Registry) ForNamespace(namespace string) *Stopwatch {
	r.mu.RLock()
	watch, ok := r.watches[namespace]
	r.mu.RUnlock()
	// Check if found.
	if ok {
		return watch
	}
	// Not found, create Stopwatch.
	r.mu.Lock()
	defer r.mu.Unlock()
	watch = &Stopwatch{
		namespace:      namespace,
		meteringPoints: make(map[string]*MeteringPoint),
	}
	r.watches[namespace] = watch
	return watch
}

// Values returns the values of all metering points.
func (r *Registry) Values() MeteringPointValues {
	mpvs := MeteringPointValues{}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, stopwatch := range r.watches {
		mpvs = append(mpvs, stopwatch.Values()...)
	}
	return mpvs
}

// Reset clears all stopwatches.
func (r *Registry) Reset() {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Simply start with a new registry, rest is done by GC.
	r.watches = make(map[string]*Stopwatch)
}

// EOF
