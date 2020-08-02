// Tideland Go Trace - Failure - Unit Tests
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package failure_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/trace/failure"
)

//--------------------
// TESTS
//--------------------

// TestNewInfoBag verifies the correct instantiation of a
// an InfoBag.
func TestNewInfoBag(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	values := map[string]interface{}{
		"a": "foo",
		"b": 12345,
		"c": 13.37,
		"d": false,
		"e": true,
	}
	ib := failure.NewInfoBag(
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
	ib := failure.NewInfoBag(
		"a", 1337,
		"b", "foo",
		"c", failure.NewInfoBag(
			"x", false,
			"y", 42,
		),
	)
	s := `[{"a": 1337}, {"b": "foo"}, {"c": [{"x": false}, {"y": 42}]}]`

	assert.Equal(ib.String(), s)
}

// EOF
