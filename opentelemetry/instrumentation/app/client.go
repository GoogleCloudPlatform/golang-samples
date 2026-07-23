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
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// callSingle makes an http request to this application's /single endpoint.
// The provided context is used to propagate the trace context with the
// http headers.
// [START opentelemetry_instrumentation_client]
func callSingle(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/single", nil)
	if err != nil {
		return err
	}
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	return res.Body.Close()
}

// [END opentelemetry_instrumentation_client]
