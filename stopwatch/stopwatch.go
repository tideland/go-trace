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
// METERING POINT
//--------------------

// MeteringPointValues contains the accumulated values of one metering point.
type MeteringPointValues struct {
	Namespace string
	ID        string
	Quantity  int64
	Total     time.Duration
	Minimum   time.Duration
	Maximum   time.Duration
	Average   time.Duration
}

// MeteringPoint collects the measurements of one code section.
type MeteringPoint struct {
	mu       sync.Mutex
	owner    *Stopwatch
	id       string
	queue    []time.Duration
	quantity int64
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

// Values returns the current values after an accumulation.
func (mp *MeteringPoint) Values() MeteringPointValues {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.accumulate()
	return MeteringPointValues{
		Namespace: mp.owner.namespace,
		ID:        mp.id,
		Quantity:  mp.quantity,
		Total:     mp.total,
		Minimum:   mp.minimum,
		Maximum:   mp.maximum,
		Average:   time.Duration(int64(mp.total) / mp.quantity),
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
	mp.reset()
}

// reset clears the measuring queue.
func (mp *MeteringPoint) reset() {
	mp.queue = make([]time.Duration, 0, 1024)
}

//--------------------
// STOPWATCHES
//--------------------

// stopwatches is the register type for all stopwatches.
type stopwatches struct {
	mu      sync.RWMutex
	watches map[string]*Stopwatch
}

// once ensures only one initialization.
var once sync.Once

// registry contains all stopwatches by ID.
var registry *stopwatches

// initializedRegistry returns the registry for the stopwatches.
func initializedRegistry() *stopwatches {
	once.Do(func() {
		if registry == nil {
			registry = &stopwatches{
				watches: make(map[string]*Stopwatch),
			}
		}
	})
	return registry
}

// load retrieves an already registered stopwatch or signals if it
// doesn't exist.
func (sws *stopwatches) load(namespace string) (*Stopwatch, bool) {
	sws.mu.RLock()
	defer sws.mu.RUnlock()
	sw, ok := sws.watches[namespace]
	return sw, ok
}

// store checks if the namespace already exists and possibly returns it.
// Otherwise it registers the given one and returns that.
func (sws *stopwatches) store(namespace string, sw *Stopwatch) *Stopwatch {
	sws.mu.Lock()
	defer sws.mu.Unlock()
	swsSW, ok := sws.watches[namespace]
	if ok {
		return swsSW
	}
	sws.watches[namespace] = sw
	return sw
}

//--------------------
// STOPWATCH
//--------------------

// Stopwatch allows to measure the execution time at multiple reading
// in one namespace.
type Stopwatch struct {
	mu             sync.RWMutex
	namespace      string
	meteringPoints map[string]*MeteringPoint
}

// New returns a new instance of a stopwatch with the given namespace.
// In case that namespace is already in use that stopwatch will be
// returned.
func New(namespace string) *Stopwatch {
	// Check for alreadoy registered stopwatch.
	sw, ok := initializedRegistry().load(namespace)
	if ok {
		return sw
	}
	// Create new stopwatch and register it.
	sw = &Stopwatch{
		namespace:      namespace,
		meteringPoints: make(map[string]*MeteringPoint),
	}
	return initializedRegistry().store(namespace, sw)
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
	mp.reset()
	sw.meteringPoints[id] = mp
	sw.mu.Unlock()
	return mp
}

// EOF
