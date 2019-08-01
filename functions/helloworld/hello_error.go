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

// Package helloworld provides a set of Cloud Functions samples.

package helloworld

import (
	"log"
	"net/http"
)

// HelloError is an HTTP Cloud Function with a request parameter.
func HelloError(w http.ResponseWriter, r *http.Request) {
	// [START functions_helloworld_error]
	// These WILL be reported to Stackdriver Error Reporting
	log.Fatal("I failed you")
	
	// These will NOT be reported to Stackdriver Error Reporting
	log.Print("I failed you")

	http.Error(w, "I failed you", 500)
	// [END functions_helloworld_error]
}

