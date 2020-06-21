// Tideland Go Trace - Crumbs - Unit Tests
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crumbs_test // import "tideland.dev/go/trace/crumbs"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/capture"
	"tideland.dev/go/trace/crumbs"
)

//--------------------
// TESTS
//--------------------

// TestNewDefault tests creating a default Crumb.
func TestNewDefault(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	c := crumbs.New()

	assert.NotNil(c)

	cw0 := c.L(0)
	cw1 := c.L(1)

	assert.NotNil(cw0)
	assert.NotNil(cw1)
	assert.Equal(cw0, cw1)
}

// TestDifferentLevelWriter creates a Crumb with level 1.
// So L() returns defferent CrumbWriter.
func TestDifferentLevelWriter(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	c := crumbs.New(crumbs.Level(1))

	cw0 := c.L(0)
	cw1 := c.L(1)
	cw2 := c.L(2)

	assert.Different(cw0, cw1)
	assert.Equal(cw1, cw2)
}

// TestDefaultWriter creates a default Crumb with the CrumWriter
// using the WriterGrainTray writing to stdout.
func TestDefaultWriter(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	cout := capture.Stdout(func() {
		c := crumbs.New()

		c.L(0).Info("info test", "a", 1, "a", 2)
	})
	assert.Contains(`"kind":"info","message":"info test","infos":[{"key":"a","value":1},{"key":"a","value":2}]`, cout.String())

	cout = capture.Stdout(func() {
		c := crumbs.New()

		c.L(0).Error(errors.New("test"), "error test", "done")
	})
	assert.Contains(`"kind":"error","message":"error test","infos":[{"key":"error","value":"test"},{"key":"done","value":true}]`, cout.String())
}

// EOF
