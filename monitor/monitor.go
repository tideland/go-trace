// Tideland Go Trace - Monitor
//
// Copyright (C) 2009-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package monitor // import "tideland.dev/go/trace/monitor"

//--------------------
// MONITOR
//--------------------

// Monitor combines StopWatch and StaySetIndicator.
type Monitor struct {
	sw  *StopWatch
	ssi *StaySetIndicator
}

// New creates a new monitor.
func New() *Monitor {
	m := &Monitor{
		sw:  newStopWatch(),
		ssi: newStaySetIndicator(),
	}
	return m
}

// StopWatch returns the internal stop watch instance.
func (m *Monitor) StopWatch() *StopWatch {
	return m.sw
}

// StaySetIndicator returns a stay-set indicator instance.
func (m *Monitor) StaySetIndicator() *StaySetIndicator {
	return m.ssi
}

// Reset clears all collected values so far.
func (m *Monitor) Reset() {
	m.sw.reset()
	m.ssi.reset()
}

// Stop terminates the monitor.
func (m *Monitor) Stop() {
	m.sw.stop()
	m.ssi.stop()
}

// EOF
