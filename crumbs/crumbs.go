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
	"os"
)

//--------------------
// CRUMBS
//--------------------

// Crumbs is the entry poing for all logging.
type Crumbs struct {
	level byte
	empty *emptyCrumbWriter
	grain *grainCrumbWriter
}

// New creates and initializes a new crumbs instances.
func New(options ...Option) *Crumbs {
	c := &Crumbs{
		level: 0,
		empty: &emptyCrumbWriter{},
		grain: &grainCrumbWriter{
			tray: NewWriterGrainTray(os.Stdout),
		},
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

// Crumble is intended to be used with defer. When called the given
// function f will be called. In case of an error an error crumb will
// be written with the given message and values.
func Crumble(cw CrumbWriter, f func() error, msg string, infos ...interface{}) {
	if err := f(); err != nil {
		if cwErr := cw.Error(err, msg, infos...); cwErr != nil {
			panic("cannot crumble error:" + cwErr.Error())
		}
	}
}

// EOF
