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
	"errors"
	"log/slog"
	"os"
)

// main runs an application which serves two http endpoints: /single and /multi.
// The calling the /multi endpoint results in multiple calls to the /single
// endpoint.
//
// The application is instrumented with OpenTelemetry and exports OTLP for
// metrics and traces. It uses log/slog for logging, and writes logs to stdout.
// The application does not include any GCP dependencies, and instead only uses
// open standards for instrumentation. The OpenTelemetry collector is used to
// route telemetry to GCP.
//
// [START opentelemetry_instrumentation_main]
func main() {
	ctx := context.Background()

	// Setup logging
	setupLogging()

	// Setup metrics, tracing, and context propagation
	shutdown, err := setupOpenTelemetry(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up OpenTelemetry", slog.Any("error", err))
		os.Exit(1)
	}

	// Run the http server, and shutdown and flush telemetry after it exits.
	slog.InfoContext(ctx, "server starting...")
	if err = errors.Join(runServer(), shutdown(ctx)); err != nil {
		slog.ErrorContext(ctx, "server exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}

// [END opentelemetry_instrumentation_main]
