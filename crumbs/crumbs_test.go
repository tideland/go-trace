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
	"bytes"
	"context"
	"errors"
	"log"
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

// TestDifferentLevelWriter tests a Crumb with level 1.
// So L() returns different CrumbWriter.
func TestDifferentLevelWriter(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	c := crumbs.New(crumbs.Level(1))

	cw0 := c.L(0)
	cw1 := c.L(1)
	cw2 := c.L(2)

	assert.Different(cw0, cw1)
	assert.Equal(cw1, cw2)
}

// TestDefaultWriter tests a default Crumb using the
// WriterGrainTray writing to stdout.
func TestDefaultWriter(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	cout := capture.Stdout(func() {
		c := crumbs.New()

		assert.NoError(c.L(0).Info("info test", "a", 1, "a", 2))
	})
	assert.Contains(`"kind":"info","message":"info test","infos":[{"key":"a","value":1},{"key":"a","value":2}]`, cout.String())

	cout = capture.Stdout(func() {
		c := crumbs.New()

		assert.NoError(c.L(0).Error(errors.New("test"), "error test", "done"))
	})
	assert.Contains(`"kind":"error","message":"error test","infos":[{"key":"error","value":"test"},{"key":"done","value":true}]`, cout.String())
}

// TestOwnWriter tests a Crumb using the WriterGrainTray
// writing to an own Writer.
func TestOwnWriter(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	buf := bytes.Buffer{}
	gt := crumbs.NewWriterGrainTray(&buf)
	c := crumbs.New(crumbs.Tray(gt))

	assert.NoError(c.L(0).Info("info test", "a", 1, "a", 2))
	assert.Contains(`"kind":"info","message":"info test","infos":[{"key":"a","value":1},{"key":"a","value":2}]`, buf.String())

	assert.NoError(c.L(0).Error(errors.New("test"), "error test", "done"))
	assert.Contains(`"kind":"error","message":"error test","infos":[{"key":"error","value":"test"},{"key":"done","value":true}]`, buf.String())
}

// TestLoggerWriter tests a Crumb using the LoggerGrainTray.
func TestLoggerWriter(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	buf := bytes.Buffer{}
	l := log.New(&buf, "crumbs", 0)
	gt := crumbs.NewLoggerGrainTray(l)
	c := crumbs.New(crumbs.Tray(gt))

	assert.NoError(c.L(0).Info("info test", "a", 1, "a", 2))
	assert.Contains(`"kind":"info","message":"info test","infos":[{"key":"a","value":1},{"key":"a","value":2}]`, buf.String())

	assert.NoError(c.L(0).Error(errors.New("test"), "error test", "done"))
	assert.Contains(`"kind":"error","message":"error test","infos":[{"key":"error","value":"test"},{"key":"done","value":true}]`, buf.String())
}

// TestContext tests the transport of a Crumb inside a Context.
func TestContext(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	cwIn := crumbs.New().L(0)
	ctxIn := context.Background()
	ctxOut := crumbs.NewContext(ctxIn, cwIn)
	cwOut, ok := crumbs.FromContext(ctxOut)

	assert.OK(ok)
	assert.Different(ctxIn, ctxOut)
	assert.Equal(cwIn, cwOut)

	assert.NoError(cwOut.Info("done"))
}

// TestCrumble tests the defer helper function Crumble.
func TestCrumble(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	buf := bytes.Buffer{}
	gt := crumbs.NewWriterGrainTray(&buf)
	cw := crumbs.New(crumbs.Tray(gt)).L(0)
	fDefer := func() {
		defer crumbs.Crumble(cw, func() error { return nil }, "ok")
		defer crumbs.Crumble(cw, func() error { return errors.New("failed") }, "failed")
	}

	fDefer()

	assert.NotContains("ok", buf.String())
	assert.Contains("failed", buf.String())
}

// EOF
