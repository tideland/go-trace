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
	"encoding/json"
	"fmt"
)

//--------------------
// INFO BAG
//--------------------

// InfoBag contains a number of useful extra informations for
// failures.
type InfoBag struct {
	infos map[string]interface{}
}

// New creates an InfoBag with a number of keys and values.
// The arguments are interpreted as a alternating pairs of those.
// So keys are taken as or converted into strings. In case the final
// item would be a key its value will be set to true.
func New(kvs ...interface{}) *InfoBag {
	ib := &InfoBag{
		infos: make(map[string]interface{}),
	}
	key := ""
	for _, kv := range kvs {
		// Check for new key.
		if key == "" {
			key = fmt.Sprintf("%s", kv)
			continue
		}
		// Now the value.
		value, ok := ib.infos[key]
		if !ok {
			// It's a new value.
			ib.infos[key] = kv
			key = ""
			continue
		}
		// Key already has one or more values.
		if values, ok := value.([]interface{}); ok {
			// Append to already known.
			ib.infos[key] = append(values, kv)
			key = ""
			continue
		}
		// So far only one value.
		ib.infos[key] = []interface{}{value, kv}
		key = ""
	}
	// Check if loop ended after key.
	if key != "" {
		ib.infos[key] = true
	}
	return ib
}

// Len returns the number of keys inside of
// the InfoBag.
func (ib InfoBag) Len() int {
	return len(ib.infos)
}

// Do iterates over the informations of the InfoBag and
// calls the given function for each key and value. The
// value will be passed as string so that in case of
// reference types they won't be modifyable.
func (ib InfoBag) Do(f func(key, value string)) {
	for k, v := range ib.infos {
		var sv string
		switch tv := v.(type) {
		case fmt.Stringer:
			sv = tv.String()
		default:
			sv = fmt.Sprintf("%v", v)
		}
		f(k, sv)
	}
}

// MarshalJSON implements the json.Marshaller interface.
// Needed for handling of nested InfoBag instances.
func (ib InfoBag) MarshalJSON() ([]byte, error) {
	return json.Marshal(ib.infos)
}

// String implements the fmt.Stringer interface.
func (ib InfoBag) String() string {
	bs, err := json.Marshal(ib.infos)
	if err != nil {
		panic("cannot stringify infobag: " + err.Error())
	}
	return string(bs)
}

// EOF
