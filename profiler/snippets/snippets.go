// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START profiler_setup_go_compute_engine]
// [START profiler_start]

// snippets is an example of starting cloud.google.com/go/profiler.
package main

import (
	"cloud.google.com/go/profiler"
)

func main() {
	// Profiler initialization, best done as early as possible.
	if err := profiler.Start(profiler.Config{
		Service:        "myservice",
		ServiceVersion: "1.0.0",
		// ProjectID must be set if not running on GCP.
		// ProjectID: "my-project",
	}); err != nil {
		// TODO: Handle error.
	}
}

// [END profiler_start]
// [END profiler_setup_go_compute_engine]
