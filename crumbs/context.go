// Tideland Go Trace - Crumbs
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crumbs // import "tideland.dev/go/trace/crumbs"

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

// crumbsContextKey is the context key for a CrumbWriter.
const crumbsContextKey contextKey = 1

// NewContext creates a context containing a CrumbWriter.
func NewContext(ctx context.Context, cw CrumbWriter) context.Context {
	return context.WithValue(ctx, crumbsContextKey, cw)
}

// FromContext retrieves a CrumbWriter from a context.
func FromContext(ctx context.Context) (CrumbWriter, bool) {
	cw, ok := ctx.Value(crumbsContextKey).(CrumbWriter)
	return cw, ok
}

// EOF
