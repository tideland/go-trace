// Tideland Go Trace - Logger - Unit Tests
//
// Copyright (C) 2012-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package logger_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/trace/logger"
)

//--------------------
// TESTS
//--------------------

// TestGetSetLevel tests the setting of the logging level.
func TestGetSetLevel(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	tw := logger.NewTestWriter()
	cw := logger.SetWriter(tw)
	defer logger.SetWriter(cw)

	logger.SetLevel(logger.LevelDebug)
	logger.Debugf("Debug.")
	logger.Infof("Info.")
	logger.Warningf("Warning.")
	logger.Errorf("Error.")
	logger.Criticalf("Critical.")

	assert.Length(tw, 5)
	tw.Reset()

	logger.SetLevel(logger.LevelError)
	logger.Debugf("Debug.")
	logger.Infof("Info.")
	logger.Warningf("Warning.")
	logger.Errorf("Error.")
	logger.Criticalf("Critical.")

	assert.Length(tw, 2)
	assert.Contents("[ERROR]", tw.Entries()[0])
	assert.Contents("[CRITICAL]", tw.Entries()[1])
	tw.Reset()
}

// TestFiltering tests the filtering of the logging.
func TestFiltering(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	tw := logger.NewTestWriter()
	cw := logger.SetWriter(tw)
	defer logger.SetWriter(cw)

	logger.SetLevel(logger.LevelDebug)
	logger.SetFilter(func(level logger.LogLevel, msg string) bool {
		return level >= logger.LevelWarning && level <= logger.LevelError
	})

	logger.Debugf("Debug.")
	logger.Infof("Info.")
	logger.Warningf("Warning.")
	logger.Errorf("Error.")
	logger.Criticalf("Critical.")

	assert.Length(tw, 2)
	tw.Reset()

	logger.UnsetFilter()

	logger.Debugf("Debug.")
	logger.Infof("Info.")
	logger.Warningf("Warning.")
	logger.Errorf("Error.")
	logger.Criticalf("Critical.")

	assert.Length(tw, 5)
	tw.Reset()
}

// TestGoLogger tests logging with the go logger.
func TestGoLogger(t *testing.T) {
	cw := logger.SetWriter(logger.NewGoWriter())
	defer logger.SetWriter(cw)

	logger.Debugf("Debug.")
	logger.Infof("Info.")
	logger.Warningf("Warning.")
	logger.Errorf("Error.")
	logger.Criticalf("Critical.")
}

// TestSysLogger tests logging with the syslogger.
func TestSysLogger(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sw, err := logger.NewSysWriter("GOTRACELOGGER")
	assert.Nil(err)
	cw := logger.SetWriter(sw)
	defer logger.SetWriter(cw)

	logger.SetLevel(logger.LevelDebug)

	logger.Debugf("Debug.")
	logger.Infof("Info.")
	logger.Warningf("Warning.")
	logger.Errorf("Error.")
	logger.Criticalf("Critical.")
}

// TestFatalExit tests the call of the fatal exiter after a
// fatal error log.
func TestFatalExit(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	tw := logger.NewTestWriter()
	cw := logger.SetWriter(tw)
	defer logger.SetWriter(cw)

	exited := false
	fatalExiter := func() {
		exited = true
	}

	logger.SetFatalExiter(fatalExiter)
	logger.Fatalf("Fatal.")

	assert.Length(tw, 1)
	assert.True(exited)
	tw.Reset()
}

// EOF
