// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_http_method]

// Package http provides a set of HTTP Cloud Function samples.
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
