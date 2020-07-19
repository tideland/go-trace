// Tideland Go Trace - Location
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package location // import "tideland.dev/go/trace/location"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

//--------------------
// LOCATION
//--------------------

// Cached locations.
var (
	mu        sync.Mutex
	locations = make(map[uintptr]Location)
)

// Location contains the formatted ID and the details
// of one location.
type Location struct {
	ID      string
	Package string
	File    string
	Func    string
	Line    int
}

// Code returns returns a location based code.
func (l Location) Code(prefix string) string {
	pparts := strings.Split(l.Package, "/")

	for _, ppart := range pparts {
		prefix += ppart[0:1]
	}

	prefix += l.File[0:1]
	prefix += strconv.Itoa(l.Line)

	return strings.ToUpper(prefix)
}

// At returns the location at the given offset.
func At(offset int) Location {
	mu.Lock()
	defer mu.Unlock()
	// Fix the offset.
	offset += 2
	if offset < 2 {
		offset = 2
	}
	// Retrieve program counters.
	pcs := make([]uintptr, 1)
	n := runtime.Callers(offset, pcs)
	if n == 0 {
		return Location{}
	}
	pcs = pcs[:n]
	// Check cache.
	pc := pcs[0]
	l, ok := locations[pc]
	if ok {
		return l
	}
	// Build ID based on program counters.
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		pkg, fun := path.Split(frame.Function)
		parts := strings.Split(fun, ".")
		pkg = path.Join(pkg, parts[0])
		fun = strings.Join(parts[1:], ".")
		_, file := path.Split(frame.File)
		id := fmt.Sprintf("(%s:%s:%s:%d)", pkg, file, fun, frame.Line)
		if !more {
			l := Location{
				ID:      id,
				Package: pkg,
				File:    file,
				Func:    fun,
				Line:    frame.Line,
			}
			locations[pc] = l
			return l
		}
	}
}

// Here return the current location.
func Here() Location {
	return At(1)
}

//--------------------
// STACK
//--------------------

// Stack contains a number of locations of a call stack.
type Stack []Location

// String returns a string representation of the stack.
func (s Stack) String() string {
	var ids []string
	for _, l := range s {
		ids = append(ids, l.ID)
	}
	return strings.Join(ids, " :: ")
}

// HereDeep returns the current callstack until the given depth.
func HereDeep(depth int) Stack {
	var stack Stack
	var end = depth + 1
	if end < 1 {
		end = 1
	}
	for i := 1; i < end; i++ {
		stack = append(stack, At(i))
	}
	return stack
}

// EOF
