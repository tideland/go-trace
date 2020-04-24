// Tideland Go Trace - Stay-set Indicator - Unit Test
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stayset_test // import "tideland.dev/go/trace/stayset"

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/trace/stayset"
)

//--------------------
// TESTS
//--------------------

// TestCreateSSI checks the creation and reusage of SSIs
// with the same ID.
func TestCreateSSI(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	ssiOne := stayset.ForNamespace("one")
	assert.NotNil(ssiOne)
	ipOneA := ssiOne.IndicatorPoint("a")
	assert.NotNil(ipOneA)

	ssiTwo := stayset.ForNamespace("two")
	assert.NotNil(ssiTwo)
	assert.Different(ssiOne, ssiTwo)
	ipTwoA := ssiTwo.IndicatorPoint("a")
	assert.Different(ipOneA, ipTwoA)

	ssiReuse := stayset.ForNamespace("one")
	assert.NotNil(ssiReuse)
	assert.Different(ssiReuse, ssiTwo)
	assert.Equal(ssiOne, ssiReuse)
}

// EOF
