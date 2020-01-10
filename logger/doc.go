// Tideland Go Trace - Logger
//
// Copyright (C) 2012-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package logger provides a flexible way to log information with different
// levels and on different backends. For tests create a test writer with
//
//     w := logger.NewTestWriter()
//     current := logger.SetWriter(w)
//     defer logger.SetWriter(current)
//
// Now logged entries can be retrieved and reseted.
//
//     es := w.Entries()
//     w.Reset()
//
// The default logger writes to stdout, others can be instantiated with
// any io.Writer. logger.NewGoWriter() returns a writer using the standard
// Go logging implementation and logger.NewSysWriter() returs a writer
// based on the system log.
//
// The levels are Debug, Info, Warning, Error, Critical, and Fatal. Here
// logger.Debugf() also logs information about file name, function
// name, and line number while log.Fatalf() may end the program
// depending on the set FatalExiterFunc.
//
// Changes to the standard behavior can be made with logger.SetLevel()
// and logger.SetFatalExiter(). Own logger backends and exiter can be
// defined. Additionally a filter function allows to drill down the
// logged entries.
package logger // import "tideland.dev/go/trace/logger"

// EOF
