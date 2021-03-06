// Tideland Go Trace - Stopwatch - Unit Test
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stopwatch_test // import "tideland.dev/go/trace/stopwatch"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/trace/stopwatch"
)

// generateMeasurings creates some measurings for the tests.
func generateMeasurings() {
	gen := generators.New(generators.FixedRand())
	swOne := stopwatch.ForNamespace("one")
	mpOneA := swOne.MeteringPoint("a")
	mpOneB := swOne.MeteringPoint("b")
	swTwo := stopwatch.ForNamespace("two")
	mpTwoA := swTwo.MeteringPoint("a")

	for i := 0; i < 777; i++ {
		m := mpOneA.Start()
		gen.SleepOneOf(1*time.Millisecond, 2*time.Millisecond)
		m.Stop()
		m = mpOneB.Start()
		gen.SleepOneOf(2*time.Millisecond, 1*time.Millisecond)
		m.Stop()
		m = mpTwoA.Start()
		gen.SleepOneOf(1*time.Millisecond, 2*time.Millisecond)
		m.Stop()
	}
}

//--------------------
// TESTS
//--------------------

// TestCreateStopwatch checks the creation and reusage of stopwatches
// with the same ID.
func TestCreateStopwatch(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	swOne := stopwatch.ForNamespace("one")
	assert.NotNil(swOne)
	mpOneA := swOne.MeteringPoint("a")
	assert.NotNil(mpOneA)

	swTwo := stopwatch.ForNamespace("two")
	assert.NotNil(swTwo)
	assert.Different(swOne, swTwo)
	mpTwoA := swTwo.MeteringPoint("a")
	assert.Different(mpOneA, mpTwoA)

	swReuse := stopwatch.ForNamespace("one")
	assert.NotNil(swReuse)
	assert.Different(swReuse, swTwo)
	assert.Equal(swOne, swReuse)
}

// TestMeasurings runs a number of measurings.
func TestMeasurings(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	generateMeasurings()

	// Only for one metering point.
	mpv := stopwatch.ForNamespace("one").MeteringPoint("a").Value()
	assert.Equal(mpv.Namespace, "one")
	assert.Equal(mpv.ID, "a")
	assert.Equal(mpv.Quantity, 777)
	assert.True(mpv.Minimum <= mpv.Average && mpv.Average <= mpv.Maximum)
	assert.Logf("%v", mpv)

	// Now for all metering points of one stopwatch.
	mpvs := stopwatch.ForNamespace("one").Values()
	assert.Length(mpvs, 2)
	for _, mpv := range mpvs {
		assert.Equal(mpv.Namespace, "one")
		assert.True(mpv.ID == "a" || mpv.ID == "b")
		assert.Equal(mpv.Quantity, 777)
		assert.True(mpv.Minimum <= mpv.Average && mpv.Average <= mpv.Maximum)
	}

	// Now for all metering points.
	mpvs = stopwatch.Values()
	assert.Length(mpvs, 3)
	for _, mpv := range mpvs {
		assert.True(mpv.Namespace == "one" || mpv.Namespace == "two")
		assert.True(mpv.ID == "a" || mpv.ID == "b")
		assert.Equal(mpv.Quantity, 777)
		assert.True(mpv.Minimum <= mpv.Average && mpv.Average <= mpv.Maximum)
	}
}

// TestReset checks the resetting of all watches.
func TestReset(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	stopwatch.Reset()
	generateMeasurings()

	// Check length.
	mpvs := stopwatch.Values()
	assert.Length(mpvs, 3)
	for _, mpv := range mpvs {
		assert.Equal(mpv.Quantity, 777)
	}

	// Reset and check length.
	stopwatch.Reset()
	mpvs = stopwatch.Values()
	assert.Length(mpvs, 0)
}

// EOF
