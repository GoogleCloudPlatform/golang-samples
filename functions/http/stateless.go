// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START functions_concepts_stateless]

// Package http provides a set of HTTP Cloud Functions samples.
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
