// Tideland Go Trace - Location
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package location provides a way to retrieve the current location in code.
// This can be used in logging or monitoring. Passing an offset helps hiding
// calling wrappers.
//
//     l1 := location.Here()
//     l2 := location.At(5)
//     id := location.At(2).ID
//     code := location.At(2).Code("ERR")
//     stack := location.HereDeep(5)
//
// Internal caching fastens retrieval after first call.
package location // import "tideland.dev/go/trace/location"

// EOF
