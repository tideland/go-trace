// Tideland Go Trace - InfoBag
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package infobag // import "tideland.dev/go/trace/infobag"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strings"
)

//--------------------
// INFO
//--------------------

// Info contains one information of the InfoBag. It consists out
// of a key and any value, which could be an InfoBag too.
type Info struct {
	Key   string
	Value interface{}
}

// String implements the fmt.Stringer interface.
func (i Info) String() string {
	switch i.Value.(type) {
	case string:
		return fmt.Sprintf("{%q: %q}", i.Key, i.Value)
	default:
		return fmt.Sprintf("{%q: %v}", i.Key, i.Value)
	}
}

//--------------------
// INFO BAG
//--------------------

// InfoBag contains a number of useful extra informations for
// failures.
type InfoBag struct {
	infos []Info
}

// New creates an InfoBag with a number of keys and values.
// The arguments are interpreted as a alternating pairs of those.
// So keys are taken as or converted into strings. In case the final
// item would be a key its value will be set to true.
func New(kvs ...interface{}) *InfoBag {
	ib := &InfoBag{}
	i := Info{}
	for _, kv := range kvs {
		// Check for new key.
		if i.Key == "" {
			i.Key = fmt.Sprintf("%s", kv)
			continue
		}
		// Now a value. Add the info to the bag.
		i.Value = kv
		ib.infos = append(ib.infos, i)
		i = Info{}
	}
	// Check if loop ended after key.
	if i.Key != "" {
		i.Value = true
		ib.infos = append(ib.infos, i)
	}
	return ib
}

// Len returns the number of informations inside of
// the InfoBag.
func (ib InfoBag) Len() int {
	return len(ib.infos)
}

// Do iterates over the informations of the InfoBag and
// calls the given function for each key and value.
func (ib InfoBag) Do(f func(key string, valiue interface{})) {
	for _, i := range ib.infos {
		f(i.Key, i.Value)
	}
}

// String implements the fmt.Stringer interface.
func (ib InfoBag) String() string {
	kvs := make([]string, len(ib.infos))
	for i, info := range ib.infos {
		kvs[i] = info.String()
	}
	kvss := strings.Join(kvs, ", ")
	return "[" + kvss + "]"
}

// EOF
