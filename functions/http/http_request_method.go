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

// [START functions_http_method]

// Package http provides a set of HTTP Cloud Functions samples.
package http

import (
	"fmt"
	"net/http"
)

// HelloHTTPMethod is an HTTP Cloud function.
// It uses the request method to differentiate the response.
func HelloHTTPMethod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Hello World!")
	case http.MethodPut:
		http.Error(w, "403 - Forbidden", http.StatusForbidden)
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// [END functions_http_method]
