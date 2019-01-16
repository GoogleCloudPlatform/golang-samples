// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"fmt"
	"net/http"
)

// [START functions_tips_scopes]

// h is in the global (instance-wide) scope.
var h string

// init runs during package initialization. So, this will only run during an
// an instance's cold start.
func init() {
	h = heavyComputation()
}

// ScopeDemo is an example of using globally and locally
// scoped variables in a function.
func ScopeDemo(w http.ResponseWriter, r *http.Request) {
	l := lightComputation()
	fmt.Fprintf(w, "Global: %q, Local: %q", h, l)
}

// [END functions_tips_scopes]

func heavyComputation() string {
	// Do something slow.
	return "slow"
}

func lightComputation() string {
	// Do something fast.
	return "fast"
}
