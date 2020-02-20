// Tideland Go Trace - Logger
//
// Copyright (C) 2012-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package logger // import "tideland.dev/go/trace/logger"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

//--------------------
// WRITER
//--------------------

// defaultTimeFormat controls how the timestamp of the standard
// logger is printed by default.
const defaultTimeFormat = "2006-01-02 15:04:05 Z07:00"

// levelText maps log levels to the according display texts.
var levelText = map[LogLevel]string{
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarning:  "WARNING",
	LevelError:    "ERROR",
	LevelCritical: "CRITICAL",
	LevelFatal:    "FATAL",
}

// levelToText translates levels to string representations.
func levelToText(level LogLevel) string {
	text, ok := levelText[level]
	if !ok {
		return "INVALID LEVEL"
	}
	return text
}

// Writer is the interface for different log writers.
type Writer interface {
	// Write writes the given message with additional
	// information at the specific log level.
	Write(level LogLevel, msg string) error
}

// standardWriter is a simple writer writing to the given I/O
// writer. Beside the output it doesn't handle the levels differently.
type standardWriter struct {
	mu         sync.Mutex
	out        io.Writer
	timeFormat string
}

// NewTimeformatWriter creates a writer writing to the passed
// output and with the specified time format.
func NewTimeformatWriter(out io.Writer, timeFormat string) Writer {
	return &standardWriter{
		out:        out,
		timeFormat: timeFormat,
	}
}

// NewStandardWriter creates the standard writer writing
// to the passed output.
func NewStandardWriter(out io.Writer) Writer {
	return NewTimeformatWriter(out, defaultTimeFormat)
}

// NewStandardOutWriter creates the standard writer writing
// to STDOUT.
func NewStandardOutWriter() Writer {
	return NewStandardWriter(os.Stdout)
}

// Write implements Writer.
func (w *standardWriter) Write(level LogLevel, msg string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().Format(w.timeFormat)
	text := levelToText(level)
	_, err := fmt.Fprintf(w.out, "%s [%s] %s\n", now, text, msg)
	return err
}

// goWriter just uses the standard go log package.
type goWriter struct{}

// NewGoWriter creates a writer using the Go log package.
func NewGoWriter() Writer {
	return &goWriter{}
}

// Write implements Writer.
func (w *goWriter) Write(level LogLevel, msg string) error {
	text := levelToText(level)
	log.Println("["+text+"]", msg)
	return nil
}

// Entries contains the collected entries of a test writer.
type Entries interface {
	// Len returns the number of collected entries.
	Len() int

	// Entries returns the collected entries.
	Entries() []string

	// Reset clears the collected entries.
	Reset()
}

// TestWriter extends the Writer interface with methods to
// retrieve and reset the collected data for testing purposes.
type TestWriter interface {
	Writer
	Entries
}

// testWriter simply collects logs to be evaluated inside of tests.
type testWriter struct {
	mu      sync.Mutex
	entries []string
}

// NewTestWriter returns a special writer for testing purposes.
func NewTestWriter() TestWriter {
	return &testWriter{}
}

// Write implements Writer.
func (w *testWriter) Write(level LogLevel, msg string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	text := levelToText(level)
	entry := fmt.Sprintf("%d [%s] %s", time.Now().UnixNano(), text, msg)
	w.entries = append(w.entries, entry)
	return nil
}

// Len implements TestWriter.
func (w *testWriter) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.entries)
}

// Entries implements TestWriter.
func (w *testWriter) Entries() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	entries := make([]string, len(w.entries))
	copy(entries, w.entries)
	return entries
}

// Reset implements TestWriter.
func (w *testWriter) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = nil
}

// EOF
