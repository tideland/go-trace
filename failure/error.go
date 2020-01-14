// Tideland Go Trace - Failure
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package failure // import "tideland.dev/go/trace/failure"

//--------------------
// IMPORTS
//--------------------

import "sync"

//--------------------
// ERROR
//--------------------

// Error encapsulates errors in a synchronized way.
type Error struct {
	mu  sync.RWMutex
	err error
}

// Set sets the encapsulated error.
func (e *Error) Set(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.err = err
}

// Get retrieves the encapsulated error.
func (e *Error) Get() error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.err
}

// IsNil checks is the encapsulated error is nil.
func (e *Error) IsNil() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.err == nil
}

// EOF
