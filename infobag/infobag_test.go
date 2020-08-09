// Tideland Go Trace - InfoBag - Unit Tests
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package infobag_test // import "tideland.dev/go/trace/infobag"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/trace/infobag"
)

//--------------------
// TESTS
//--------------------

// TestNew verifies the correct instantiation of a
// an InfoBag.
func TestNew(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	values := map[string]interface{}{
		"a": "foo",
		"b": 12345,
		"c": 13.37,
		"d": false,
		"e": true,
	}
	ib := infobag.New(
		"a", values["a"],
		"b", values["b"],
		"c", values["c"],
		"d", values["d"],
		"e",
	)

	assert.Length(ib, 5)

	ib.Do(func(key string, value interface{}) {
		assert.Equal(value, values[key])
	})
}

// TestInfoBagString verifies the string output of
// InfoBags.
func TestInfoBagString(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ib := infobag.New(
		"a", 1337,
		"b", "foo",
		"c", infobag.New(
			"x", false,
			"y", 42,
		),
		"b", "bar",
	)
	s := `{"a":1337,"b":["foo","bar"],"c":{"x":false,"y":42}}`

	// TODO Better test, map order may change.
	assert.Equal(ib.String(), s)
}

// EOF
