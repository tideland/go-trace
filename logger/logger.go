// Tideland Go Trace - Logger
//
// Copyright (C) 2012-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package logger // import "tideland.dev/go/trace/logger"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"os"
	"sync"

	"tideland.dev/go/trace/location"
)

//--------------------
// LEVEL
//--------------------

// LogLevel describes the chosen log level between
// debug and critical.
type LogLevel int

// Log levels to control the logging output.
const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
	LevelFatal
)

//--------------------
// EXIT
//--------------------

// FatalExiterFunc defines a functions that will be called
// in case of a Fatalf call.
type FatalExiterFunc func()

// OSFatalExiter exits the application with os.Exit and
// the return code -1.
func OSFatalExiter() {
	os.Exit(-1)
}

// PanicFatalExiter exits the application with a panic.
func PanicFatalExiter() {
	panic("program aborted after fatal situation, see log")
}

//--------------------
// FILTER
//--------------------

// FilterFunc allows to filter the output of the logging. Filters
// have to return true if the received entry shall be filtered and
// not output.
type FilterFunc func(level LogLevel, msg string) bool

//--------------------
// LOGGER API
//--------------------

// Level returns the current log level.
func Level() LogLevel {
	backend.mu.RLock()
	defer backend.mu.RUnlock()
	return backend.level
}

// SetLevel sets the log level to a new one and returns the current.
func SetLevel(level LogLevel) LogLevel {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	current := backend.level
	switch {
	case level <= LevelDebug:
		backend.level = LevelDebug
	case level >= LevelFatal:
		backend.level = LevelFatal
	default:
		backend.level = level
	}
	return current
}

// SetWriter sets the writing target to a new one and returns the current.
func SetWriter(out Writer) Writer {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	current := backend.out
	if out != nil {
		backend.out = out
	}
	return current
}

// SetFatalExiter sets the fatal exiter function to a new one and returns the current.
func SetFatalExiter(fef FatalExiterFunc) FatalExiterFunc {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	current := backend.fatalExiter
	if fef != nil {
		backend.fatalExiter = fef
	}
	return current
}

// SetFilter sets the global output filter to a new one and returns the current.
// Nil function is allowed, it unsets the filter.
func SetFilter(ff FilterFunc) FilterFunc {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	current := backend.shallWrite
	backend.shallWrite = ff
	return current
}

// UnsetFilter removes the global output filter and returns the current.
func UnsetFilter() FilterFunc {
	return SetFilter(nil)
}

// Debugf logs a message at debug level.
func Debugf(format string, args ...interface{}) {
	backend.log(LevelDebug, location.At(1).ID+" "+format, args...)
}

// Infof logs a message at info level.
func Infof(format string, args ...interface{}) {
	backend.log(LevelInfo, format, args...)
}

// Warningf logs a message at warning level.
func Warningf(format string, args ...interface{}) {
	backend.log(LevelWarning, format, args...)
}

// Errorf logs a message at error level.
func Errorf(format string, args ...interface{}) {
	backend.log(LevelError, format, args...)
}

// Criticalf logs a message at critical level.
func Criticalf(format string, args ...interface{}) {
	backend.log(LevelCritical, location.At(1).ID+" "+format, args...)
}

// Fatalf logs a message at fatal level. After logging the message the
// function calls the fatal exiter function, which by default means exiting
// the application with error code -1. So only call in real fatal cases.
func Fatalf(format string, args ...interface{}) {
	backend.log(LevelFatal, location.At(1).ID+" "+format, args...)
	backend.mu.Lock()
	defer backend.mu.Unlock()
	backend.fatalExiter()
}

//--------------------
// LOGGER IMPLEMENTATION
//--------------------

// loggerBackend provides a flexible configurable logging system.
type loggerBackend struct {
	mu          sync.RWMutex
	level       LogLevel
	out         Writer
	fatalExiter FatalExiterFunc
	shallWrite  FilterFunc
}

// log checks level and filter and performs the logging.
func (lb *loggerBackend) log(level LogLevel, format string, args ...interface{}) {
	// Copies to not block the logger.
	lb.mu.RLock()
	lbLevel := lb.level
	lbShallWrite := lb.shallWrite
	lb.mu.RUnlock()
	if lbLevel > level {
		// Passed level is too low.
		return
	}
	msg := fmt.Sprintf(format, args...)
	if lbShallWrite != nil && !lbShallWrite(level, msg) {
		// Filter rejects log entry.
		return
	}
	lb.mu.Lock()
	_ = lb.out.Write(level, msg)
	lb.mu.Unlock()
}

// backend provides the logger backend. It is initialised with
// info level, using stdout for writing, and ends with os.Exit(-1)
// in case of a fatal entry.
var backend = &loggerBackend{
	level:       LevelInfo,
	out:         NewStandardOutWriter(),
	fatalExiter: OSFatalExiter,
}

// EOF
