// Tideland Go Trace - Monitor
//
// Copyright (C) 2009-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package monitor // import "tideland.dev/go/trace/monitor"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"sync"
	"time"

	"tideland.dev/go/trace/failure"
)

//--------------------
// INDICATOR VALUE
//--------------------

// IndicatorValue manages the value range for one indicator.
type IndicatorValue struct {
	ID      string
	Count   int
	Current int
	Min     int
	Max     int
}

// String implements fmt.Stringer.
func (iv *IndicatorValue) String() string {
	return fmt.Sprintf("%s: %d / act %d / min %d / max %d", iv.ID, iv.Count, iv.Current, iv.Min, iv.Max)
}

// update the indicator value.
func (iv *IndicatorValue) update(shallIncr bool) {
	// Check for initial values.
	if iv.Count == 0 {
		iv.Count = 1
		iv.Current = 1
		iv.Min = 1
		iv.Max = 1
	}
	// Regular update.
	iv.Count++
	if shallIncr {
		iv.Current++
	} else {
		iv.Current--
	}
	if iv.Current < iv.Min {
		iv.Min = iv.Current
	}
	if iv.Current > iv.Max {
		iv.Max = iv.Current
	}
}

// IndicatorValues is a set of stay-set values.
type IndicatorValues []IndicatorValue

// Implement the sort interface.

func (ivs IndicatorValues) Len() int           { return len(ivs) }
func (ivs IndicatorValues) Swap(i, j int)      { ivs[i], ivs[j] = ivs[j], ivs[i] }
func (ivs IndicatorValues) Less(i, j int) bool { return ivs[i].ID < ivs[j].ID }

//--------------------
// STAY-SET INDICATOR
//--------------------

// Describing increment or decrement of stay-set values.
const (
	up   = true
	down = false
)

// StaySetIndicator allows to increase and decrease stay-set values.
type StaySetIndicator struct {
	actionC chan func()
	doneC   chan struct{}
	changes map[string][]bool
	values  map[string]*IndicatorValue
}

// newStaySetIndicator creates a new StaySetIndicator.
func newStaySetIndicator() *StaySetIndicator {
	i := &StaySetIndicator{
		actionC: make(chan func(), 128),
		doneC:   make(chan struct{}),
		changes: make(map[string][]bool),
		values:  make(map[string]*IndicatorValue),
	}
	go i.backend()
	return i
}

// Increase increases a stay-set staySetIndicator.
func (i *StaySetIndicator) Increase(id string) {
	i.actionC <- func() {
		i.changes[id] = append(i.changes[id], up)
	}
}

// Decrease decreases a stay-set staySetIndicator.
func (i *StaySetIndicator) Decrease(id string) {
	i.actionC <- func() {
		i.changes[id] = append(i.changes[id], down)
	}
}

// Read returns a stay-set staySetIndicator.
func (i *StaySetIndicator) Read(id string) (IndicatorValue, error) {
	var iv *IndicatorValue
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	i.actionC <- func() {
		defer wg.Done()
		i.accumulateOne(id)
		iv = i.values[id]
		if iv == nil {
			err = failure.New("indicator value '%s' does not exist", id)
		}
	}
	wg.Wait()
	if iv == nil {
		return IndicatorValue{}, err
	}
	return *iv, nil
}

// Do performs the function f for all values.
func (i *StaySetIndicator) Do(f func(IndicatorValue) error) error {
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	i.actionC <- func() {
		defer wg.Done()
		i.accumulateAll()
		for _, ssi := range i.values {
			if err = f(*ssi); err != nil {
				return
			}
		}
	}
	wg.Wait()
	return err
}

// reset clears all values.
func (i *StaySetIndicator) reset() {
	i.actionC <- func() {
		i.changes = make(map[string][]bool)
		i.values = make(map[string]*IndicatorValue)
	}
}

// stop terminates the indicator.
func (i *StaySetIndicator) stop() {
	close(i.doneC)
}

// backend rund the indicator backend goroutine.
func (i *StaySetIndicator) backend() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-i.doneC:
			return
		case action := <-i.actionC:
			action()
		case <-ticker.C:
			i.accumulateAll()
		}
	}
}

// accumulateOne updates the indicator value for one ID.
func (i *StaySetIndicator) accumulateOne(id string) {
	changes, ok := i.changes[id]
	if ok {
		iv := i.values[id]
		if iv == nil {
			iv = &IndicatorValue{
				ID: id,
			}
			i.values[id] = iv
		}
		for _, increment := range changes {
			iv.update(increment)
		}
		i.changes[id] = []bool{}
	}
}

// accumulateAll updates all indicator values.
func (i *StaySetIndicator) accumulateAll() {
	for id := range i.changes {
		i.accumulateOne(id)
	}
}

// EOF
