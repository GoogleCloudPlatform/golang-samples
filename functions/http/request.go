// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_concepts_requests]
// [START functions_tips_connection_pooling]

// Package http provides a set of HTTP Cloud Function samples.
package http

import (
	"fmt"
	"net/http"
	"time"
)

var urlString = "https://example.com"

// client is used to make HTTP requests with a 10 second timeout.
// http.Clients should be reused instead of created as needed.
var client = &http.Client{
	Timeout: 10 * time.Second,
}

// MakeRequest is an example of making an HTTP request. MakeRequest uses a
// single http.Client for all requests to take advantage of connection
// pooling and caching. See https://godoc.org/net/http#Client.
func MakeRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := client.Get(urlString)
	if err != nil {
		http.Error(w, "Error making request", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Bad StatusCode: %d", resp.StatusCode)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ok")
}

// [END functions_tips_connection_pooling]
// [END functions_concepts_requests]
