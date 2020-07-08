// Tideland Go Trace - Stay-set Indicator
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stayset // import "tideland.dev/go/trace/stayset"

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

// registryContextKey is the context key for a Registry.
const registryContextKey contextKey = 1

// NewContext creates a context containing a Registry.
func NewContext(ctx context.Context, r *Registry) context.Context {
	return context.WithValue(ctx, registryContextKey, r)
}

// FromContext retrieves a Registry from a context.
func FromContext(ctx context.Context) (*Registry, bool) {
	r, ok := ctx.Value(registryContextKey).(*Registry)
	return r, ok
}

// EOF
