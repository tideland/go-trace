// Tideland Go Trace - Stay-set Indicator - Unit Test
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stayset_test // import "tideland.dev/go/trace/stayset"

import (
	"context"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/trace/stayset"
)

// generateIndicators creates some indicators for the tests.
func generateIndicators(r *stayset.Registry) [3]int {
	gen := generators.New(generators.FixedRand())
	ssiOne := r.ForNamespace("one")
	ssiOneA := ssiOne.IndicatorPoint("a")
	ssiOneB := ssiOne.IndicatorPointWithValue("b", 10)
	ssiTwo := r.ForNamespace("two")
	ssiTwoA := ssiTwo.IndicatorPoint("a")
	points := []*stayset.IndicatorPoint{ssiOneA, ssiOneB, ssiTwoA}
	pointQueues := [3][]stayset.Indication{}
	quantities := [3]int{}

	for j := 0; j < 2500; j++ {
		b := gen.OneByteOf(0, 1, 2, 1, 2, 2)

		if gen.FlipCoin(40) {
			// Start indication.
			p := points[b]
			i := p.Start()
			pointQueues[b] = append(pointQueues[b], i)
			quantities[b]++

			continue
		}
		// Stop indication.
		if len(pointQueues[b]) > 0 {
			i := pointQueues[b][0]
			pointQueues[b] = pointQueues[b][1:]

			i.Stop()
		}
	}

	return quantities
}

//--------------------
// TESTS
//--------------------

// TestCreateSSI checks the creation and reusage of SSIs
// with the same ID.
func TestCreateSSI(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	r := stayset.New()

	ssiOne := r.ForNamespace("one")
	assert.NotNil(ssiOne)
	ipOneA := ssiOne.IndicatorPoint("a")
	assert.NotNil(ipOneA)

	ssiTwo := r.ForNamespace("two")
	assert.NotNil(ssiTwo)
	assert.Different(ssiOne, ssiTwo)
	ipTwoA := ssiTwo.IndicatorPoint("a")
	assert.Different(ipOneA, ipTwoA)

	ssiReuse := r.ForNamespace("one")
	assert.NotNil(ssiReuse)
	assert.Different(ssiReuse, ssiTwo)
	assert.Equal(ssiOne, ssiReuse)
}

// TestIndicators runs a number of SSIs.
func TestIndicators(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	r := stayset.New()
	quantities := generateIndicators(r)

	// Only for one indicator.
	ipv := r.ForNamespace("one").IndicatorPoint("a").Value()
	assert.Equal(ipv.Namespace, "one")
	assert.Equal(ipv.ID, "a")
	assert.Range(ipv.Quantity, 0, quantities[0])
	assert.Equal(ipv.Minimum, 0)
	assert.Range(ipv.Maximum, 0, ipv.Quantity)
	assert.Logf("%v", ipv)

	// Only for one indicator with pre-set value.
	ipv = r.ForNamespace("one").IndicatorPoint("b").Value()
	assert.Equal(ipv.Namespace, "one")
	assert.Equal(ipv.ID, "b")
	assert.Range(ipv.Quantity, 0, quantities[1])
	assert.Equal(ipv.Minimum, 10)
	assert.Range(ipv.Maximum, 10, ipv.Quantity)
	assert.Logf("%v", ipv)

	// Now for all indicators of one namespace.
	ipvs := r.ForNamespace("one").Values()

	assert.Length(ipvs, 2)

	for _, ipv := range ipvs {
		assert.Equal(ipv.Namespace, "one")
		assert.True(ipv.ID == "a" || ipv.ID == "b")
	}

	// Now for all indicators points.
	ipvs = r.Values()

	assert.Length(ipvs, 3)

	for _, ipv := range ipvs {
		assert.True(ipv.Namespace == "one" || ipv.Namespace == "two")
		assert.True(ipv.ID == "a" || ipv.ID == "b")
	}
}

// TestReset checks the resetting of all SSIs.
func TestReset(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	r := stayset.New()
	_ = generateIndicators(r)

	// Check length.
	ipvs := r.Values()
	assert.Length(ipvs, 3)

	// Reset and check length.
	r.Reset()
	ipvs = r.Values()
	assert.Length(ipvs, 0)
}

// TestContext tests the transport of a SSI inside a Context.
func TestContext(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	rIn := stayset.New()
	ctxIn := context.Background()
	ctxOut := stayset.NewContext(ctxIn, rIn)
	rOut, ok := stayset.FromContext(ctxOut)

	assert.OK(ok)
	assert.Different(ctxIn, ctxOut)
	assert.Equal(rIn, rOut)
	assert.NotNil(rOut)
}

// EOF
