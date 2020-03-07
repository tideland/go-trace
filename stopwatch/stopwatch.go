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

//--------------------
// METERING POINT
//--------------------

// MeteringPoint collects the measurements of one code section.
type MeteringPoint struct {
	mu         sync.RWMutex
	id         string
	queueIndex int
	queue      []time.Duration
	quantity   int
	total      time.Duration
	minimum    time.Duration
	maximum    time.Duration
}

// Start begins a new measurement.
func (mp MeteringPoint) Start() Measurement {
	return Measurement{
		owner: &mp,
		start: time.Now(),
	}
}

//--------------------
// STOPWATCH
//--------------------

// Stopwatch allows to measure the execution time at multiple reading
// in one namespace.
type Stopwatch struct {
	mu             sync.RWMutex
	namespace      string
	meteringPoints map[string]MeteringPoint
}

// New returns a new instance of a stopwatch with the given namespace.
// In case that namespace is already in use that stopwatch will be
// returned.
func New(namespace string) *Stopwatch {
	// Check for alreadoy registered stopwatch.
	// TODO
	// Create new stopwatch and register it.
	sw := &Stopwatch{
		namespace:      namespace,
		meteringPoints: make(map[string]MeteringPoint),
	}
	return sw
}

// MeteringPoint returns a new or already existing metering point
// with the given ID.
func (sw *Stopwatch) MeteringPoint(id string) MeteringPoint {
	// First check for existing metering point.
	sw.mu.RLock()
	mp, ok := sw.meteringPoints[id]
	sw.mu.RUnlock()
	if ok {
		return mp
	}
	// Not yet existing.
	sw.mu.Lock()
	mp = MeteringPoint{
		id:    id,
		queue: make([]time.Duration, 1024),
	}
	sw.meteringPoints[id] = mp
	sw.mu.Unlock()
	return mp
}

// EOF
