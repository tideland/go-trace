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

	"tideland.dev/go/audit/asserts"
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

// EOF
