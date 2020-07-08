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
// STAY-SET INDICATOR
//--------------------

// SSI allows to measure the execution numbers at multiple counting points
// in one namespace.
type SSI struct {
	mu              sync.RWMutex
	namespace       string
	indicatorPoints map[string]*IndicatorPoint
}

// IndicatorPointWithValue returns a new or already existing indicator point
// with the given ID and and initial value. In case it's already existing the
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

//--------------------
// REGISTRY
//--------------------

// Registry is the register type for a number of SSIs.
type Registry struct {
	mu   sync.RWMutex
	ssis map[string]*SSI
}

// New returns a registry for stay-set indicators.
func New() *Registry {
	r := &Registry{
		ssis: make(map[string]*SSI),
	}
	return r
}

// ForNamespace returns an instance of a SSI with the given namespace.
// In case that namespace is already in use that SSI will be returned.
func (r *Registry) ForNamespace(namespace string) *SSI {
	r.mu.RLock()
	ssi, ok := r.ssis[namespace]
	r.mu.RUnlock()
	// Check if found.
	if ok {
		return ssi
	}
	// Not found, create SSI.
	r.mu.Lock()
	defer r.mu.Unlock()
	ssi = &SSI{
		namespace:       namespace,
		indicatorPoints: make(map[string]*IndicatorPoint),
	}
	r.ssis[namespace] = ssi
	return ssi
}

// Values returns the values of all metering points.
func (r *Registry) Values() IndicatorPointValues {
	ipvs := IndicatorPointValues{}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ssi := range r.ssis {
		ipvs = append(ipvs, ssi.Values()...)
	}
	return ipvs
}

// Reset clears all namespaces.
func (r *Registry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ssis = make(map[string]*SSI)
}

// EOF
