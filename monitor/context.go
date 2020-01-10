// Tideland Go Trace - Monitor
//
// Copyright (C) 2009-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package monitor // import "tideland.dev/go/trace/monitor"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
)

//--------------------
// CONTEXT
//--------------------

// contextKey describes the type of the context key.
type contextKey int

// monitorContextKey is the context key for a logger.
const monitorContextKey contextKey = 1

// NewContext creates a context containing a logger.
func NewContext(ctx context.Context, monitor Monitor) context.Context {
	return context.WithValue(ctx, monitorContextKey, monitor)
}

// FromContext retrieves a logger from a context.
func FromContext(ctx context.Context) (Monitor, bool) {
	monitor, ok := ctx.Value(monitorContextKey).(Monitor)
	return monitor, ok
}

// EOF
