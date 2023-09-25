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
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func handleSingle(w http.ResponseWriter, r *http.Request) {
	sleepTime := time.Duration(100+rand.Intn(100)) * time.Millisecond
	time.Sleep(sleepTime)
	fmt.Fprintf(w, "slept %v\n", sleepTime)
}

func handleMulti(w http.ResponseWriter, r *http.Request) {
	subRequests := 3 + rand.Intn(4)
	slog.InfoContext(r.Context(), "handle /multi request", slog.Int("subRequests", subRequests))

	for i := 0; i < subRequests; i++ {
		if err := callSingle(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}

	fmt.Fprintln(w, "ok")
}

func runServer() error {
	handleHTTP("/single", handleSingle)
	handleHTTP("/multi", handleMulti)
	return http.ListenAndServe(":8080", nil)
}

func handleHTTP(route string, handleFn http.HandlerFunc) {
	http.Handle(route, otelhttp.NewHandler(otelhttp.WithRouteTag(route, handleFn), route))
}
