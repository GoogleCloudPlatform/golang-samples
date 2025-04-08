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

// [START profiler_setup_go_compute_engine]

// snippets is an example of starting cloud.google.com/go/profiler.
package main

import (
	"cloud.google.com/go/profiler"
)

func main() {
	cfg := profiler.Config{
		Service:        "myservice",
		ServiceVersion: "1.0.0",
		// ProjectID must be set if not running on GCP.
		// ProjectID: "my-project",

		// For OpenCensus users:
		// To see Profiler agent spans in APM backend,
		// set EnableOCTelemetry to true
		// EnableOCTelemetry: true,
	}

	// Profiler initialization, best done as early as possible.
	if err := profiler.Start(cfg); err != nil {
		// TODO: Handle error.
	}
}

// [END profiler_setup_go_compute_engine]
