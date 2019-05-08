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

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"fmt"
	"net/http"
)

// [START functions_tips_scopes]
// [START run_tips_global_scope]

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

// [END run_tips_global_scope]
// [END functions_tips_scopes]

func heavyComputation() string {
	// Do something slow.
	return "slow"
}

func lightComputation() string {
	// Do something fast.
	return "fast"
}
