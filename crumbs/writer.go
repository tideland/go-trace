// Tideland Go Trace - Crumbs
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crumbs // import "tideland.dev/go/trace/crumbs"

//--------------------
// CRUMB WRITER
//--------------------

// CrumbWriter is a leveled and configured writer for infos
// and errors.
type CrumbWriter interface {
	// Info logs messages with key/value pairs of additional
	// information.
	Info(msg string, keysAndValues ...interface{})

	// Error logs errors with an additional message and key/value
	// pairs of additional information.
	Error(err error, msg string, keysAndValues ...interface{})
}

//--------------------
// EMPTY CRUMP WRITER
//--------------------

// emptyCrumbWriter simply does nothing as its (virtual) level is
// lower than the crumb level.
type emptyCrumbWriter struct{}

// Info implements CrumbWriter.
func (w *emptyCrumbWriter) Info(msg string, keysAndValues ...interface{}) {}

// Error implements CrumbWriter.
func (w *emptyCrumbWriter) Error(err error, msg string, keysAndValues ...interface{}) {}

//--------------------
// GRAIN CRUMB WRITER
//--------------------

// grainCrumbWriter creates the Grains for the configurable backend.
type grainCrumbWriter struct {
}

// Info implements CrumbWriter.
func (w *grainCrumbWriter) Info(msg string, keysAndValues ...interface{}) {}

// Error implements CrumbWriter.
func (w *grainCrumbWriter) Error(err error, msg string, keysAndValues ...interface{}) {}

// EOF
