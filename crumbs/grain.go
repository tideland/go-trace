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
	"time"
)

//--------------------
// GRAIN KIND
//--------------------

// GrainKind describes if a Grain is an information or an error.
type GrainKind int

// Kind of grains.
const (
	InfoGrain GrainKind = iota
	ErrorGrain
)

// String implements fmt.Stringer.
func (gk GrainKind) String() string {
	switch gk {
	case InfoGrain:
		return "info"
	case ErrorGrain:
		return "error"
	default:
		return "unknown"
	}
}

// MarshalJSON implements json.Marshaler.
func (gk GrainKind) MarshalJSON() ([]byte, error) {
	return []byte(`"` + gk.String() + `"`), nil
}

//--------------------
// GRAIN
//--------------------

// GrainInfo contains a pair of key and value of a Grain.
type GrainInfo struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Grain contains all data to log.
type Grain struct {
	Timestamp time.Time   `json:"timestamp"`
	Kind      GrainKind   `json:"kind"`
	Message   string      `json:"message"`
	Infos     []GrainInfo `json:"infos"`
}

// newGrain parses the keys and values and creates a Grain.
func newGrain(kind GrainKind, msg string, infos ...interface{}) *Grain {
	g := &Grain{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Message:   msg,
	}
	key := ""
	last := len(infos) - 1
	for i, kv := range infos {
		switch {
		case i%2 == 0 && i == last:
			g.Infos = append(g.Infos, GrainInfo{
				Key:   fmt.Sprintf("%v", kv),
				Value: true,
			})
		case i%2 == 0:
			key = fmt.Sprintf("%v", kv)
		default:
			g.Infos = append(g.Infos, GrainInfo{
				Key:   key,
				Value: kv,
			})
			key = ""
		}
	}
	return g
}

// String implements fmt.Stringer. This implementation
// marshals the Grain into JSON.
func (g Grain) String() string {
	b, err := json.Marshal(g)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// EOF
