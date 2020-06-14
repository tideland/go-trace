// Tideland Go Trace - Crumbs
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crumbs // import "tideland.dev/go/trace/crumbs"

// Option is a function able to configure the Crumbs.
type Option func(c *Crumbs)

// Level sets the Crumbs level for logging. Default is 0.
func Level(level byte) Option {
	return func(c *Crumbs) {
		if level >= 0 {
			c.level = level
		}
	}
}

// Tray sets the GrainTray backend. Default is the StandardLoggerTray.
func Tray(tray GrainTray) Option {
	return func(c *Crumbs) {
		if tray != nil {
			c.grain.tray = tray
		}
	}
}

// EOF
