// Tideland Go Trace - Failure
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package failure // import "tideland.dev/go/trace/failure"

//--------------------
// IMPORTS
//--------------------

//--------------------
// INFO
//--------------------

// info contains one information of the InfoBag. It consists out
// of a key and any value, which could be an InfoBag too.
type info struct {
	key   string
	value interface{}
}

//--------------------
// INFO BAG
//--------------------

// EOF
