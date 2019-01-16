// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_concepts_stateless]

// Package http provides a set of HTTP Cloud Function samples.
package http

import (
	"fmt"
	"net/http"
)

// count is a global variable, but only shared within a function instance.
var count = 0

// ExecutionCount is an HTTP Cloud Function that counts how many times it
// is executed within a specific instance.
func ExecutionCount(w http.ResponseWriter, r *http.Request) {
	count++

	// Note: the total function invocation count across
	// all instances may not be equal to this value!
	fmt.Fprintf(w, "Instance execution count: %d", count)
}

// [END functions_concepts_stateless]
