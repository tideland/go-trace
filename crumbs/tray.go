// Tideland Go Trace - Crumbs
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crumbs // import "tideland.dev/go/trace/crumbs"

//--------------------
// GRAIN TRAY
//--------------------

// GrainTray defines the interface for the backends responsible
// to collect and store the Grains.
type GrainTray interface {
	// Put adds a Grain to the tray.
	Put(grain *Grain) error
}

// EOF
