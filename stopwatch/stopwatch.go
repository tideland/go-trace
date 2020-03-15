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

// MeteringPoint collects the measurements of one code section.
type MeteringPoint struct {
	mu         sync.RWMutex
	owner      *Stopwatch
	id         string
	queueIndex int
	queue      []time.Duration
	quantity   int
	total      time.Duration
	minimum    time.Duration
	maximum    time.Duration
}

// Start begins a new measurement.
func (mp *MeteringPoint) Start() Measurement {
	return Measurement{
		owner: mp,
		start: time.Now(),
	}
}

// enqueue adds the measured duration.
func (mp *MeteringPoint) enqueue(measuring time.Duration) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	if mp.queueIndex == cap(mp.queue) {
		// Accumulate the enqueued values.
		measurings := mp.queue
		mp.queue = make([]time.Duration, 1024)
		mp.queueIndex = 0
		go mp.accumulate(measurings, nil)
	}
	mp.queue[mp.queueIndex] = measuring
	mp.queueIndex++
}

// accumulateNow synchronously evaluates the measurings.
func (mp *MeteringPoint) accumulateNow() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	var wg sync.WaitGroup
	wg.Add(1)
	measurings := mp.queue
	mp.queue = make([]time.Duration, 1024)
	mp.queueIndex = 0
	go mp.accumulate(measurings, &wg)
	wg.Wait()
}

// accumulate evaluates the collected measurings synchronous or asynchronous.
func (mp *MeteringPoint) accumulate(measurings []time.Duration, wg *sync.WaitGroup) {
	// Get initial values.
	mp.mu.RLock()
	quantity := mp.quantity
	total := mp.total
	minimum := mp.minimum
	maximum := mp.maximum
	mp.mu.RUnlock()
	// Accumulate set of measurings isolated.
	for _, duration := range measurings {
		if quantity == 0 {
			quantity = 1
			total = duration
			minimum = duration
			maximum = duration
			continue
		}
		quantity++
		total += duration
		if minimum > duration {
			minimum = duration
		}
		if maximum < duration {
			maximum = duration
		}
	}
	// Total update.
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.quantity += quantity
	mp.total += total
	if mp.minimum > minimum {
		mp.minimum = minimum
	}
	if mp.maximum < maximum {
		mp.maximum = maximum
	}
	// Tell waiting caller that it's done.
	if wg != nil {
		wg.Done()
	}
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
		queue: make([]time.Duration, 1024),
	}
	sw.meteringPoints[id] = mp
	sw.mu.Unlock()
	return mp
}

// EOF
