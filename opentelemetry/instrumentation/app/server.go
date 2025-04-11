// Copyright 2023 Google LLC
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

package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// handleSingle handles an http request by sleeping for 100-200 ms. It writes
// the number of milliseconds slept as its response.
// [START opentelemetry_instrumentation_handle_single]
func handleSingle(w http.ResponseWriter, r *http.Request) {
	sleepTime := randomSleep(r)
	fmt.Fprintf(w, "work completed in %v\n", sleepTime)
}

// [END opentelemetry_instrumentation_handle_single]

// handleMulti handles an http request by making 3-7 http requests to the
// /single endpoint.
// [START opentelemetry_instrumentation_handle_multi]
func handleMulti(w http.ResponseWriter, r *http.Request) {
	subRequests := 3 + rand.Intn(4)
	// Write a structured log with the request context, which allows the log to
	// be linked with the trace for this request.
	slog.InfoContext(r.Context(), "handle /multi request", slog.Int("subRequests", subRequests))

	err := computeSubrequests(r, subRequests)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	fmt.Fprintln(w, "ok")
}

// [END opentelemetry_instrumentation_handle_multi]

// runServer runs an http server on port 8080 which handles requests to the
// /multi and /single endpoints.
// [START opentelemetry_instrumentation_run_server]
func runServer() error {
	handleHTTP("/single", handleSingle)
	handleHTTP("/multi", handleMulti)

	return http.ListenAndServe(":8080", nil)
}

// handleHTTP handles the http HandlerFunc on the specified route, and uses
// otelhttp for context propagation, trace instrumentation, and metric
// instrumentation.
func handleHTTP(route string, handleFn http.HandlerFunc) {
	instrumentedHandler := otelhttp.NewHandler(otelhttp.WithRouteTag(route, handleFn), route)

	http.Handle(route, instrumentedHandler)
}

// [END opentelemetry_instrumentation_run_server]
