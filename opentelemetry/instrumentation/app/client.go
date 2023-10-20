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

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// callSingle makes an http request to this application's /single endpoint.
// The provided context is used to propagate the trace context with the
// http headers.
// [START opentelemetry_instrumentation_client]
func callSingle(ctx context.Context) error {
	// otelhttp.Get makes an http GET request, just like net/http.Get.
	// In addition, it records a span, records metrics, and propagates context.
	res, err := otelhttp.Get(ctx, "http://localhost:8080/single")
	if err != nil {
		return err
	}

	return res.Body.Close()
}

// [END opentelemetry_instrumentation_client]
