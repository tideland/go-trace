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
// MEASURINGS
//--------------------

// Measuring defines one execution time measuring containing the ID and
// the starting time of the measuring and able to pass this data after
// the end of the measuring to the measurer.
type Measuring struct {
	owner *StopWatch
	id    string
	begin time.Time
}

// End ends the measuring and passes it to the measurer.
func (m *Measuring) End() time.Duration {
	duration := time.Since(m.begin)
	m.owner.end(m.id, duration)
	return duration
}

// WatchValue manages the value range for one watch.
type WatchValue struct {
	ID    string
	Count int
	Total time.Duration
	Min   time.Duration
	Max   time.Duration
	Avg   time.Duration
}

// String implements fmt.Stringer.
func (wv *WatchValue) String() string {
	factor := 1000000.0
	total := float64(wv.Total) / factor
	min := float64(wv.Min) / factor
	max := float64(wv.Max) / factor
	avg := float64(wv.Avg) / factor
	return fmt.Sprintf("%s: %d / total %.4f ms / min %.4f ms / max %.4f ms / avg %.4f ms", wv.ID, wv.Count, total, min, max, avg)
}

// update the value with a new measured duration.
func (wv *WatchValue) update(duration time.Duration) {
	// Check for initial values.
	if wv.Count == 0 {
		wv.Count = 1
		wv.Total = duration
		wv.Min = duration
		wv.Max = duration
		wv.Avg = duration
		return
	}
	// Regular update.
	wv.Count++
	wv.Total += duration
	if wv.Min > duration {
		wv.Min = duration
	}
	if wv.Max < duration {
		wv.Max = duration
	}
	wv.Avg = time.Duration(int64(wv.Total) / int64(wv.Count))
}

// WatchValues is a set of values.
type WatchValues []WatchValue

// Implement the sort interface.

func (wvs WatchValues) Len() int           { return len(wvs) }
func (wvs WatchValues) Swap(i, j int)      { wvs[i], wvs[j] = wvs[j], wvs[i] }
func (wvs WatchValues) Less(i, j int) bool { return wvs[i].ID < wvs[j].ID }

//--------------------
// STOP WATCH
//--------------------

// StopWatch allows to measure the execution time of
// code fragments.
type StopWatch struct {
	actionC    chan func()
	doneC      chan struct{}
	measurings map[string][]time.Duration
	values     map[string]*WatchValue
}

// newStopWatch creates a new stop watch.
func newStopWatch() *StopWatch {
	s := &StopWatch{
		actionC:    make(chan func(), 128),
		doneC:      make(chan struct{}),
		measurings: make(map[string][]time.Duration),
		values:     make(map[string]*WatchValue),
	}
	go s.backend()
	return s
}

// Begin starts a new measuring with a given id.
func (s *StopWatch) Begin(id string) *Measuring {
	return &Measuring{
		owner: s,
		id:    id,
		begin: time.Now(),
	}
}

// end returns a measuring to the collected ones.
func (s *StopWatch) end(id string, duration time.Duration) {
	s.actionC <- func() {
		s.measurings[id] = append(s.measurings[id], duration)
	}
}

// Measure measures the execution time of one function.
func (s *StopWatch) Measure(id string, f func()) time.Duration {
	m := s.Begin(id)
	f()
	return m.End()
}

// Read returns the measuring point for an id.
func (s *StopWatch) Read(id string) (WatchValue, error) {
	var wv *WatchValue
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	s.actionC <- func() {
		defer wg.Done()
		s.accumulateOne(id)
		wv = s.values[id]
		if wv == nil {
			err = failure.New("watch value '%s' does not exist", id)
		}
	}
	wg.Wait()
	if wv == nil {
		return WatchValue{}, err
	}
	return *wv, nil
}

// Do performs the function f for all measuring points.
func (s *StopWatch) Do(f func(WatchValue) error) error {
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	s.actionC <- func() {
		defer wg.Done()
		s.accumulateAll()
		for _, wv := range s.values {
			if err = f(*wv); err != nil {
				return
			}
		}
	}
	wg.Wait()
	return err
}

// reset clears all values.
func (s *StopWatch) reset() {
	s.actionC <- func() {
		s.measurings = make(map[string][]time.Duration)
		s.values = make(map[string]*WatchValue)
	}
}

// stop terminates the indicator.
func (s *StopWatch) stop() {
	close(s.doneC)
}

// backend rund the indicator backend goroutine.
func (s *StopWatch) backend() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-s.doneC:
			return
		case action := <-s.actionC:
			action()
		case <-ticker.C:
			s.accumulateAll()
		}
	}
}

// accumulateOne updates one watch value.
func (s *StopWatch) accumulateOne(id string) {
	measurings, ok := s.measurings[id]
	if ok {
		wv := s.values[id]
		if wv == nil {
			wv = &WatchValue{
				ID: id,
			}
			s.values[id] = wv
		}
		for _, duration := range measurings {
			wv.update(duration)
		}
		s.measurings[id] = []time.Duration{}
	}
}

// accumulateAll accumulates all watch values.
func (s *StopWatch) accumulateAll() {
	for id := range s.measurings {
		s.accumulateOne(id)
	}
}

// EOF
