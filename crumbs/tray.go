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
	"encoding/json"
	"fmt"
	"io"
	"time"
)

//--------------------
// GRAIN TRAY
//--------------------

// GrainTray defines the interface for the backends responsible
// to collect and store the Grains.
type GrainTray interface {
	// Put adds a Grain to the tray.
	Put(grain *Grain) error
}

//--------------------
// WRITER GRAIN TRAY
//--------------------

// WriterGrainTray writes the grain to the configured writer.
type WriterGrainTray struct {
	out io.Writer
}

// NewWriterGrainTray creates a WriterGrainTray with the given writer.
func NewWriterGrainTray(out io.Writer) *WriterGrainTray {
	return &WriterGrainTray{
		out: out,
	}
}

// Put implements GrainTray.
func (t *WriterGrainTray) Put(grain *Grain) error {
	ts := grain.Timestamp.Format(time.RFC3339Nano)
	ks := []string{"info", "error"}[grain.Kind]
	kvs, err := json.Marshal(grain.KeyValues)
	if err != nil {
		return fmt.Errorf("writer grain tray: cannot put grain: %v", err)
	}
	_, err = fmt.Fprintf(t.out, "%s (%s) %s %s\n", ts, ks, grain.Message, string(kvs))
	if err != nil {
		return fmt.Errorf("writer grain tray: cannot put grain: %v", err)
	}
	return nil
}

// EOF
