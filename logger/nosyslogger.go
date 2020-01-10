// Tideland Go Trace - Logger - No SysLogger
//
// Copyright (C) 2012-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// +build windows plan9 nacl

package logger // import "tideland.dev/go/trace/logger"

//--------------------
// IMPORTS
//--------------------

import (
	"log"
)

//--------------------
// SYSWRITER
//--------------------

// sysWriter uses the Go syslog package. It does not work
// on Windows or Plan9.
type nosyslogWriter struct {
	tag string
}

// NewSysWriter creates a writer using the Go syslog package.
// It does not work on Windows or Plan9. Here the Go log
// package is used.
func NewSysWriter(tag string) (Writer, error) {
	if len(tag) > 0 {
		tag = "(" + tag + ")"
	}
	return &nosyslogWriter{
		tag: tag,
	}, nil
}

// Write implements Writer.
func (w *nosyslogWriter) Write(level LogLevel, msg string) {
	text := levelToText(level)
	log.Println("["+text+"]", msg)
}

// EOF
