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
	"errors"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/trace/failure"
)

//--------------------
// TESTS
//--------------------

// TestIsError tests the creation and checking of errors.
func TestIsError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	emsg := "test error %d"
	err := failure.New(emsg, 1)

	assert.True(failure.IsValid(err))
	assert.Equal(err.Error(), "[ETGTFF31] test error 1")

	err = testError("test error 2")

	assert.False(failure.IsValid(err))
	assert.ErrorMatch(err, "test error 2")

	err = errors.New("42")

	assert.False(failure.IsValid(err))
}

// TestValidation checks the validation of errors and
// the retrieval of details.
func TestValidation(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// First a valid error.
	emsg := "valid"
	err := failure.New(emsg)
	assert.True(failure.IsValid(err))

	hereID, lerr := failure.Location(err)
	assert.Nil(lerr)
	assert.Equal(hereID, "(tideland.dev/go/trace/failure_test:failure_test.go:TestValidation:53)")

	// Now an invalid error.
	err = errors.New("ouch")
	assert.False(failure.IsValid(err))

	hereID, lerr = failure.Location(err)
	assert.Equal(lerr.Error(), "[ETGTFF156] passed error has invalid type: ouch")
	assert.Empty(hereID)
}

// TestAnnotation the annotation of errors with new errors.
func TestAnnotation(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	err1 := testError("wrapped")
	err2 := failure.Annotate(err1, "1st annotated")
	err3 := failure.Annotate(err2, "2nd annotated")

	assert.ErrorMatch(err3, `.* 2nd annotated: .* 1st annotated: wrapped`)
	assert.Equal(failure.Annotated(err3), err2)
	assert.Equal(failure.Annotated(err2), err1)
	assert.Length(failure.Stack(err3), 3)

	err4 := failure.Annotate(nil, "not existing")
	assert.NoError(err4)
}

// TestFirst tests choosing the first (existing) error.
func TestFirst(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	errA := testError("one")
	errB := testError("two")

	err := failure.First(nil, errA)
	assert.ErrorMatch(err, "one")
	err = failure.First(err, errB)
	assert.ErrorMatch(err, "one")
}

// TestCollection tests the collection of multiple errors to one.
func TestCollection(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	errA := testError("one")
	errB := testError("two")
	errC := testError("three")
	errD := testError("four")
	cerr := failure.Collect(errA, errB, errC, errD)

	assert.ErrorMatch(cerr, "one :: two :: three :: four")

	cerr = failure.Collect(errA, errB, nil, errD)

	assert.ErrorMatch(cerr, "one :: two :: four")

	cerr = failure.Collect()

	assert.NoError(cerr)
}

// TestDoAll tests the iteration over errors.
func TestDoAll(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msgs := []string{}
	f := func(err error) {
		msgs = append(msgs, err.Error())
	}

	// Test it on annotated errors.
	errX := testError("init")
	errA := failure.Annotate(errX, "foo")
	errB := failure.Annotate(errA, "bar")
	errC := failure.Annotate(errB, "baz")
	errD := failure.Annotate(errC, "yadda")

	failure.DoAll(errD, f)

	assert.Length(msgs, 5)

	// Test it on collected errors.
	msgs = []string{}
	errA = testError("foo")
	errB = testError("bar")
	errC = testError("baz")
	errD = testError("yadda")
	cerr := failure.Collect(errA, errB, errC, errD)

	failure.DoAll(cerr, f)

	assert.Equal(msgs, []string{"foo", "bar", "baz", "yadda"})

	// Test it on a single error.
	msgs = []string{}
	errA = testError("foo")

	failure.DoAll(errA, f)

	assert.Equal(msgs, []string{"foo"})
}

//--------------------
// HELPERS
//--------------------

type testError string

func (e testError) Error() string {
	return string(e)
}

// EOF
