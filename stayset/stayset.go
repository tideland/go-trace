// Tideland Go Trace - Stay-set Indicator
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stayset // import "tideland.dev/go/trace/stayset"

//--------------------
// IMPORT
//--------------------

import (
	"fmt"
	"sync"
)

//--------------------
// INDICATION
//--------------------

// Indication contains one current indication point.
type Indication struct {
	owner *IndicatorPoint
}

// Stop ends the counting.
func (i *Indication) Stop() {
	i.owner.stopIndication()
}

//--------------------
// INDICATOR POINT VALUE
//--------------------

// IndicatorPointValue contains the accumulated value of one indicator point.
type IndicatorPointValue struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id"`
	Quantity  int    `json:"quantity"`
	Current   int    `json:"current"`
	Minimum   int    `json:"minimum"`
	Maximum   int    `json:"maximum"`
}

// String implements the fmt.Stringer interface.
func (ipv IndicatorPointValue) String() string {
	return fmt.Sprintf(
		"[%s :: %s] %d / %v / %v / %v",
		ipv.Namespace,
		ipv.ID,
		ipv.Quantity,
		ipv.Current,
		ipv.Minimum,
		ipv.Maximum,
	)
}

// IndicatorPointValues contains a set of accumulated indicator point values.
type IndicatorPointValues []IndicatorPointValue

//--------------------
// INDICATOR POINT
//--------------------

// IndicatorPoint collects the countings of one code section.
type IndicatorPoint struct {
	mu       sync.Mutex
	owner    *SSI
	id       string
	quantity int
	current  int
	minimum  int
	maximum  int
}

// Start begins a new measurement.
func (ip *IndicatorPoint) Start() Indication {
	ip.mu.Lock()
	defer ip.mu.Unlock()
	ip.quantity++
	ip.current++
	if ip.current > ip.maximum {
		ip.maximum = ip.current
	}
	return Indication{
		owner: ip,
	}
}

// Measure is a convenience function to measure one anonymous function.
func (ip *IndicatorPoint) Measure(f func()) {
	i := ip.Start()
	defer i.Stop()
	f()
}

// Value returns the current value.
func (ip *IndicatorPoint) Value() IndicatorPointValue {
	ip.mu.Lock()
	defer ip.mu.Unlock()
	return IndicatorPointValue{
		Namespace: ip.owner.namespace,
		ID:        ip.id,
		Quantity:  ip.quantity,
		Current:   ip.current,
		Minimum:   ip.minimum,
		Maximum:   ip.maximum,
	}
}

// stopIndication ends the indication and reduces the values.
func (ip *IndicatorPoint) stopIndication() {
	ip.mu.Lock()
	defer ip.mu.Unlock()
	ip.current--
	if ip.current < ip.minimum {
		ip.minimum = ip.current
	}
}

//--------------------
// REGISTRY
//--------------------

// ssis is the register type for all SSIs.
type ssis struct {
	mu   sync.RWMutex
	ssis map[string]*SSI
}

// once ensures only one initialization.
var once sync.Once

// registry contains all indicators by namespace.
var registry *ssis

// initializedRegistry returns the registry for the stopwatches.
func initializedRegistry() *ssis {
	once.Do(func() {
		if registry == nil {
			registry = &ssis{
				ssis: make(map[string]*SSI),
			}
		}
	})
	return registry
}

// load retrieves an already registered indicator or signals if it
// doesn't exist.
func (s *ssis) load(namespace string) (*SSI, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.ssis[namespace]
	return i, ok
}

// store checks if the namespace already exists and possibly returns it.
// Otherwise it registers the given one and returns that.
func (s *ssis) store(namespace string, ssi *SSI) *SSI {
	s.mu.Lock()
	defer s.mu.Unlock()
	issi, ok := s.ssis[namespace]
	if ok {
		return issi
	}
	s.ssis[namespace] = ssi
	return ssi
}

// Values returns the values of all metering points.
func Values() IndicatorPointValues {
	ipvs := IndicatorPointValues{}
	reg := initializedRegistry()
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	for _, ssi := range reg.ssis {
		ipvs = append(ipvs, ssi.Values()...)
	}
	return ipvs
}

// Reset clears all indicators.
func Reset() {
	reg := initializedRegistry()
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	// Simply start with a new registry, rest is done by GC.
	reg.ssis = make(map[string]*SSI)
}

//--------------------
// STAY-SET INDICATOR
//--------------------

// SSI allows to measure the execution numbers at multiple counting points
// in one namespace.
type SSI struct {
	mu              sync.RWMutex
	namespace       string
	indicatorPoints map[string]*IndicatorPoint
}

// ForNamespace returns an instance of a SSI with the given namespace.
// In case that namespace is already in use that SSI will be returned.
func ForNamespace(namespace string) *SSI {
	// Check for alreadoy registered stopwatch.
	ssi, ok := initializedRegistry().load(namespace)
	if ok {
		return ssi
	}
	// Create new SSI and register it.
	ssi = &SSI{
		namespace:       namespace,
		indicatorPoints: make(map[string]*IndicatorPoint),
	}
	return initializedRegistry().store(namespace, ssi)
}

// IndicatorPointWithValue returns a new or already existing indicator point
// with the given ID and and intial value. In case it's already existing the
// value isn't changed anymore.
func (ssi *SSI) IndicatorPointWithValue(id string, value int) *IndicatorPoint {
	// First check for existing metering point.
	ssi.mu.RLock()
	ip, ok := ssi.indicatorPoints[id]
	ssi.mu.RUnlock()
	if ok {
		return ip
	}
	// Not yet existing.
	ssi.mu.Lock()
	ip = &IndicatorPoint{
		owner:   ssi,
		id:      id,
		current: value,
		minimum: value,
		maximum: value,
	}
	ssi.indicatorPoints[id] = ip
	ssi.mu.Unlock()
	return ip
}

// IndicatorPoint returns a new or already existing indicator point with
// the given ID.
func (ssi *SSI) IndicatorPoint(id string) *IndicatorPoint {
	return ssi.IndicatorPointWithValue(id, 0)
}

// Values returns the accumulated counting point values of this SSI.
func (ssi *SSI) Values() IndicatorPointValues {
	ssi.mu.RLock()
	defer ssi.mu.RUnlock()
	ipvs := make(IndicatorPointValues, len(ssi.indicatorPoints))
	i := 0
	for _, ip := range ssi.indicatorPoints {
		ipvs[i] = ip.Value()
		i++
	}
	return ipvs
}

// EOF
