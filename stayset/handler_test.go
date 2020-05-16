// Tideland Go Trace - Stay-set Indicator - Unit Test
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package stayset_test // import "tideland.dev/go/trace/stayset"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/environments"
	"tideland.dev/go/trace/stayset"
)

//--------------------
// WEB ASSERTER
//--------------------

// StartTestServer initialises and starts the asserter for the tests.
func startWebAsserter(assert *asserts.Asserts) *environments.WebAsserter {
	wa := environments.NewWebAsserter(assert)
	return wa
}

//--------------------
// TESTS
//--------------------

// TestWebValues tests retrieving the values via web handler.
func TestWebValues(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	_ = generateIndicators()

	wa.Handle("/stayset/", stayset.NewHandler())

	wreq := wa.CreateRequest(http.MethodGet, "/stayset/")
	wresp := wreq.Do()
	wresp.AssertStatusCodeEquals(http.StatusOK)
	wresp.Header().AssertKeyContainsValue("Content-Type", environments.ContentTypeJSON)
	wresp.AssertBodyContains(`"namespace":"one"`)
	wresp.AssertBodyContains(`"namespace":"two"`)
	wresp.AssertBodyContains(`"id":"a"`)
}

// TestWebReset tests resetting the values via web handler.
func TestWebReset(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	_ = generateIndicators()

	wa.Handle("/stayset/", stayset.NewHandler())

	wreq := wa.CreateRequest(http.MethodDelete, "/stayset/")
	wresp := wreq.Do()
	wresp.AssertStatusCodeEquals(http.StatusOK)
	wresp.Header().AssertKeyContainsValue("Content-Type", environments.ContentTypeJSON)
	wresp.AssertBodyContains(`"indicator point values resetted"`)
}

// TestWebIllegal tests the handler with an illegal HTTP method.
func TestWebIllegal(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	wa.Handle("/stayset/", stayset.NewHandler())

	wreq := wa.CreateRequest(http.MethodPost, "/stayset/")
	wresp := wreq.Do()
	wresp.AssertStatusCodeEquals(http.StatusMethodNotAllowed)
	wresp.AssertBodyContains("only GET and DELETE allowed")
}

// EOF
