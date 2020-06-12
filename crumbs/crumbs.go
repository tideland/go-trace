// Tideland Go Trace - Crumbs
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crumbs // import "tideland.dev/go/trace/crumbs"

//--------------------
// CRUMBS
//--------------------

// Crumbs is the entry poing for all logging.
type Crumbs struct {
	level byte
	empty CrumbsWriter
	grain CrumbsWriter
}

// New creates and initializes a new crumbs instances.
func New(options ...Option) *Crumbs {
	c := &Crumbs{
		level: 0,
		empty: &emptyCrumbWriter{},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// L returns a leveled crumb writer.
func (c *Crumbs) L(level byte) CrumbWriter {
	if level < c.level {
		return c.empty
	}
	return c.grain
}

// EOF