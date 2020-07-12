// Tideland Go Trace - Stopwatch
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stopwatch // import "tideland.dev/go/trace/stopwatch"

//--------------------
// IMPORT
//--------------------

import (
	"encoding/json"
	"net/http"
)

//--------------------
// HANDLER
//--------------------

// Handler implements the http.Handler.
type Handler struct {
	r *Registry
}

// NewHandler returns an instance of a web handler for the
// stopwatch.
func NewHandler(r *Registry) Handler {
	return Handler{
		r: r,
	}
}

// ServeHTTP implements the handling function.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Retrieve all metering point values.
		mpvs := h.r.Values()
		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		err := enc.Encode(mpvs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		// Reset all metering point values.
		h.r.Reset()
		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		err := enc.Encode("metering point values resetted")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "only GET and DELETE allowed", http.StatusMethodNotAllowed)
}

// EOF
