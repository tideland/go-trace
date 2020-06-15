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
	"fmt"
	"time"
)

//--------------------
// CONSTANTS
//--------------------

// GrainKind describes if a Grain is an information or an error.
type GrainKind int

// Kind of grains.
const (
	InfoGrain GrainKind = iota
	ErrorGrain
)

//--------------------
// GRAIN
//--------------------

// GrainKeyValue contains one of the key/values pairs or a Grain.
type GrainKeyValue struct {
	Key   string
	Value interface{}
}

// Grain contains all data to log.
type Grain struct {
	Timestamp time.Time
	Kind      GrainKind
	Message   string
	KeyValues []GrainKeyValue
}

// newGrain parses the keys and values and creates a Grain.
func newGrain(kind GrainKind, msg string, keysAndValues ...interface{}) *Grain {
	g := &Grain{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Message:   msg,
	}
	key := ""
	last := len(keysAndValues) - 1
	for i, kv := range keysAndValues {
		switch {
		case i%2 == 0 && i == last:
			g.KeyValues = append(g.KeyValues, GrainKeyValue{
				Key:   fmt.Sprintf("%v", kv),
				value: true,
			})
		case i%2 == 0:
			key = fmt.Sprintf("%v", kv)
		default:
			g.KeyValues = append(g.KeyValues, GrainKeyValue{
				Key:   key,
				value: kv,
			})
			key = ""
		}
	}
	return g
}

// EOF
