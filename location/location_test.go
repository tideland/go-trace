// Tideland Go Trace - Location - Unit Tests
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package location_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/trace/location"
)

//--------------------
// TESTS
//--------------------

// TestAt tests retrieving the location in a detailed
// way and as ID.
func TestAt(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	l := location.At(0)

	assert.Equal(l.Package, "tideland.dev/go/trace/location_test")
	assert.Equal(l.File, "location_test.go")
	assert.Equal(l.Func, "TestAt")
	assert.Equal(l.Line, 30)

	id := location.At(0).ID

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:TestAt:37)")

	code := location.At(0).Code("ERR:")

	assert.Equal(code, "ERR:TGTLL41")
}

// TestHere tests retrieving the location in a detailed
// way and as ID.
func TestHere(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	l := location.Here()

	assert.Equal(l.Package, "tideland.dev/go/trace/location_test")
	assert.Equal(l.File, "location_test.go")
	assert.Equal(l.Func, "TestHere")
	assert.Equal(l.Line, 51)

	id := location.Here().ID

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:TestHere:58)")

	code := location.Here().Code("ERR:")

	assert.Equal(code, "ERR:TGTLL62")
}

// TestStack tests retrieving a call stack.
func TestStack(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	s := stackOne()

	assert.Length(s, 5)
	assert.Equal(s[0].ID, "(tideland.dev/go/trace/location_test:location_test.go:stackFive:160)")
	assert.Equal(s[1].ID, "(tideland.dev/go/trace/location_test:location_test.go:stackFour:155)")
	assert.Equal(s[2].ID, "(tideland.dev/go/trace/location_test:location_test.go:stackThree:150)")
	assert.Equal(s[3].ID, "(tideland.dev/go/trace/location_test:location_test.go:stackTwo:145)")
	assert.Equal(s[4].ID, "(tideland.dev/go/trace/location_test:location_test.go:stackOne:140)")
}

// TestOffset tests retrieving the location with an offset.
func TestOffset(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	id := there()

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:TestOffset:85)")

	id = nestedThere()

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:TestOffset:89)")

	id = nameless()

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:nameless.func1:133)")

	id = location.At(-5).ID

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:TestOffset:97)")
}

// TestCache tests the caching of locations.
func TestCache(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	for i := 0; i < 100; i++ {
		id := nameless()

		assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:nameless.func1:133)")
	}
}

//--------------------
// HELPER
//--------------------

// there returns the id at the location of the caller.
func there() string {
	return location.At(1).ID
}

// nestedThere returns the id at the location of the caller but inside a local func.
func nestedThere() string {
	where := func() string {
		return location.At(2).ID
	}
	return where()
}

// nameless returns the id from calling a nested nameless function w/o an offset.
func nameless() string {
	noname := func() string {
		return location.Here().ID
	}
	return noname()
}

// stackOne is the first one of a stack calling function set.
func stackOne() location.Stack {
	return stackTwo()
}

// stackTwo is the second one of a stack calling function set.
func stackTwo() location.Stack {
	return stackThree()
}

// stackThree is the third one of a stack calling function set.
func stackThree() location.Stack {
	return stackFour()
}

// stackFour is the fourth one of a stack calling function set.
func stackFour() location.Stack {
	return stackFive()
}

// stackFive is the fifth one of a stack calling function set.
func stackFive() location.Stack {
	return location.HereDeep(5)
}

// EOF
