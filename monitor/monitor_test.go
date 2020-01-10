// Tideland Go Trace - Monitor - Unit Tests
//
// Copyright (C) 2009-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package monitor_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/trace/monitor"
)

//--------------------
// TESTS
//--------------------

// TestSimpleMonitor test creating and stopping a monitor.
func TestSimpleMonitor(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	m := monitor.New()
	defer m.Stop()

	assert.True(m.StopWatch().Measure("simple", func() { time.Sleep(time.Millisecond) }) > 0)

	mp, err := m.StopWatch().Read("simple")
	assert.NoError(err)
	assert.Equal(mp.ID, "simple")
}

// TestStopWatch tests the stop watch.
func TestStopWatch(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	gen := generators.New(generators.FixedRand())
	m := monitor.New()
	defer m.Stop()

	// Generate some measurings first.
	for i := 0; i < 500; i++ {
		m.StopWatch().Measure("watch", func() {
			gen.SleepOneOf(1*time.Millisecond, 3*time.Millisecond, 10*time.Millisecond)
		})
	}

	sw, err := m.StopWatch().Read("doesnotexist")
	assert.ErrorMatch(err, `.* watch value 'doesnotexist' does not exist`)

	// Check access of one measuring point.
	sw, err = m.StopWatch().Read("watch")
	assert.Nil(err)
	assert.Equal(sw.ID, "watch")
	assert.Equal(sw.Count, 500)
	assert.True(sw.Min <= sw.Avg && sw.Avg <= sw.Max)

	// Check iteration over all measuring points.
	wvs := monitor.WatchValues{}
	err = m.StopWatch().Do(func(wv monitor.WatchValue) error {
		wvs = append(wvs, wv)
		return nil
	})
	assert.Nil(err)
	assert.Length(wvs, 1)

	// Check resetting the measurings.
	m.Reset()

	wvs = monitor.WatchValues{}
	err = m.StopWatch().Do(func(wv monitor.WatchValue) error {
		wvs = append(wvs, wv)
		return nil
	})
	assert.Nil(err)
	assert.Empty(wvs)
}

// Test of the stay-set indicators  of the monitor.
func TestStaySetIndicators(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	gen := generators.New(generators.FixedRand())
	m := monitor.New()
	defer m.Stop()

	// Generate some indicators first.
	for i := 0; i < 500; i++ {
		id := gen.OneStringOf("foo", "bar", "baz", "yadda", "deadbeef")
		if gen.FlipCoin(60) {
			m.StaySetIndicator().Increase(id)
		} else {
			m.StaySetIndicator().Decrease(id)
		}
	}

	iv, err := m.StaySetIndicator().Read("doesnotexist")
	assert.ErrorMatch(err, `.* indicator value 'doesnotexist' does not exist`)

	// Check access of one stay-set indicator.
	iv, err = m.StaySetIndicator().Read("foo")
	assert.Nil(err)
	assert.Equal(iv.ID, "foo")
	assert.Equal(iv.Count, 99)
	assert.True(iv.Min <= iv.Current && iv.Current <= iv.Max)

	// Check iteration over all measuring points.
	ivs := monitor.IndicatorValues{}
	err = m.StaySetIndicator().Do(func(iv monitor.IndicatorValue) error {
		ivs = append(ivs, iv)
		return nil
	})
	assert.Nil(err)
	assert.Length(ivs, 5)

	// Check resetting the measurings.
	m.Reset()

	ivs = monitor.IndicatorValues{}
	err = m.StaySetIndicator().Do(func(iv monitor.IndicatorValue) error {
		ivs = append(ivs, iv)
		return nil
	})
	assert.Nil(err)
	assert.Empty(ivs)
}

//--------------------
// BENCHMARKS
//--------------------

// BenchmarkMonitor checks the performance of monitor.
func BenchmarkMonitor(b *testing.B) {
	gen := generators.New(generators.SimpleRand())
	m := monitor.New()
	defer m.Stop()

	for i := 0; i < b.N; i++ {
		m.StopWatch().Measure("bench", func() {
			gen.SleepOneOf(1*time.Millisecond, 3*time.Millisecond, 5*time.Millisecond)
		})
	}
}

// EOF
