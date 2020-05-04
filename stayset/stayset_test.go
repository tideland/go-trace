// Tideland Go Trace - Stay-set Indicator - Unit Test
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stayset_test // import "tideland.dev/go/trace/stayset"

import (
	"sync"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/trace/stayset"
)

// generateIndicators creates some indicators for the tests.
func generateIndicators() {
	gen := generators.New(generators.FixedRand())
	ssiOne := stayset.ForNamespace("one")
	ssiOneA := ssiOne.IndicatorPoint("a")
	ssiOneB := ssiOne.IndicatorPointWithValue("b", 10)
	ssiTwo := stayset.ForNamespace("two")
	ssiTwoA := ssiTwo.IndicatorPoint("a")
	points := []*stayset.IndicatorPoint{ssiOneA, ssiOneB, ssiTwoA}

	var wg sync.WaitGroup
	wg.Add(2500)

	for j := 0; j < 2500; j++ {
		go func() {
			b := gen.OneByteOf(0, 1, 2, 1, 2, 2)
			point := points[b]
			i := point.Start()
			gen.SleepOneOf(1*time.Millisecond, 2*time.Millisecond, 4*time.Millisecond)
			i.Stop()
			wg.Done()
		}()
	}

	wg.Wait()
}

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

// TestIndicators runs a number of stay-set indications.
func TestIndicators(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	generateIndicators()

	// Only for one indicator.
	iv := stayset.ForNamespace("one").IndicatorPoint("a").Value()
	assert.Equal(iv.Namespace, "one")
	assert.Equal(iv.ID, "a")
	assert.Range(iv.Quantity, 400, 450)
	assert.Equal(iv.Minimum, 0)
	assert.Range(iv.Maximum, 250, 300)
	assert.Logf("%v", iv)
}

// EOF
