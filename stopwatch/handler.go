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
	"fmt"
	"net/http"
)

//--------------------
// HANDLER
//--------------------

// HandlerFunc implements the net/http
func HandlerFunc(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Retrieve all metering point values.
		mpvs := Values()
		enc := json.NewEncoder(rw)
		err := enc.Encode(mpvs)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		return
	case http.MethodDelete:
		// Reset all metering point values.
		Reset()
		fmt.Fprintf(rw, "metering point values resetted")
		return
	}
	http.Error(rw, "only GET and DELETE allowed", http.StatusMethodNotAllowed)
}

// EOF
