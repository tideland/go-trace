// Tideland Go Trace - Stopwatch
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package stopwatch helps measuring execution times of Go applications. It
// measures the times of wanted code sections and calculates minimum, maximum,
// average and total values.
//
// Individual stopwatches can be created via New and a namespace. They are
// automatically registered at a backend for the global retrieval. The
// measuring points got individual identifiers.
package stopwatch // import "tideland.dev/go/trace/stopwatch"

// EOF
