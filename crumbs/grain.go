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
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Grain contains all data to log.
type Grain struct {
	Timestamp time.Time       `json:"timestamp"`
	Kind      GrainKind       `json:"kind"`
	Message   string          `json:"message"`
	KeyValues []GrainKeyValue `json:"key_values"`
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
				Value: true,
			})
		case i%2 == 0:
			key = fmt.Sprintf("%v", kv)
		default:
			g.KeyValues = append(g.KeyValues, GrainKeyValue{
				Key:   key,
				Value: kv,
			})
			key = ""
		}
	}
	return g
}

// EOF
