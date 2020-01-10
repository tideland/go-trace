// Tideland Go Trace - Failure
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package failure allows to create more detailed error values than
// with errors.New() or fmt.Errorf().
//
// The errors package allows to easily created formatted errors
// with New() like with the fmt.Errorf() function, but also containing
// the location and a code. The creation of error nessages is like with
// fmt.Errorf().
//
// If an error alreay exists use Annotate(). This way the original
// error will be stored and can be retrieved with Annotated(). Also
// its error message will be appended to the created error separated
// by a colon.
//
// All errors additionally contain their package, filename and line
// number. These information can be retrieved using Location(). In
// case of a chain of annotated errors those can be retrieved as a
// slice of errors with Stack().
package failure // import "tideland.dev/go/trace/failure"

// EOF
