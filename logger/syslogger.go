// Tideland Go Trace - Logger - Writer using syslog
//
// Copyright (C) 2012-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// +build !windows,!nacl,!plan9

package logger // import "tideland.dev/go/trace/logger"

//--------------------
// IMPORTS
//--------------------

import (
	"log"
	"log/syslog"
)

//--------------------
// SYSWRITER
//--------------------

// syslogWriter uses the Go syslog package. It does not work
// on Windows or Plan9.
type syslogWriter struct {
	out *syslog.Writer
}

// NewSysWriter creates a writer using the Go syslog package.
// It does not work on Windows or Plan9. Here the Go log
// package is used.
func NewSysWriter(tag string) (Writer, error) {
	out, err := syslog.New(syslog.LOG_DEBUG|syslog.LOG_LOCAL0, tag)
	if err != nil {
		log.Fatalf("cannot init syslog: %v", err)
		return nil, err
	}
	return &syslogWriter{
		out: out,
	}, nil
}

// Write implements Writer.
func (w *syslogWriter) Write(level LogLevel, msg string) {
	switch level {
	case LevelDebug:
		w.out.Debug(msg)
	case LevelInfo:
		w.out.Info(msg)
	case LevelWarning:
		w.out.Warning(msg)
	case LevelError:
		w.out.Err(msg)
	case LevelCritical:
		w.out.Crit(msg)
	case LevelFatal:
		w.out.Emerg(msg)
	default:
		w.out.Warning("[INVALID]" + msg)
	}
}

// EOF
