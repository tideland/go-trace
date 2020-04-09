// Tideland Go Trace - Stopwatch - Unit Test
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stopwatch_test

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

//--------------------
// TESTS
//--------------------

// TestCreateStopwatch checks the creation and reusage of stopwatches
// with the same ID.
func TestCreateStopwatch(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	swOne := stopwatch.New("one")
	assert.NotNil(swOne)
	mpOneA := swOne.MeteringPoint("a")
	assert.NotNil(mpOneA)

	swTwo := stopwatch.New("two")
	assert.NotNil(swTwo)
	assert.Different(swOne, swTwo)
	mpTwoA := swTwo.MeteringPoint("a")
	assert.Different(mpOneA, mpTwoA)

	swReuse := stopwatch.New("one")
	assert.NotNil(swReuse)
	assert.Different(swReuse, swTwo)
	assert.Equal(swOne, swReuse)
}

// TestMeasurings runs a number of measurings.
func TestMeasurings(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	gen := generators.New(generators.FixedRand())

	swOne := stopwatch.New("one")
	mpOneA := swOne.MeteringPoint("a")

	for i := 0; i < 2500; i++ {
		m := mpOneA.Start()
		gen.SleepOneOf(1*time.Millisecond, 3*time.Millisecond, 10*time.Millisecond)
		m.Stop()
	}

	values := mpOneA.Values()
	assert.Equal(values.Namespace, "one")
	assert.Equal(values.ID, "a")
	assert.Equal(values.Quantity, int64(2500))
	assert.True(values.Minimum <= values.Average && values.Average <= values.Maximum)
}

// EOF
