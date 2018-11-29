// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START profiler_setup_go_appengine]

// appengine is an example of starting cloud.google.com/go/profiler on
// App Engine.
package main

import (
	"cloud.google.com/go/profiler"
)

func main() {
	// Profiler initialization, best done as early as possible.
	if err := profiler.Start(profiler.Config{
		// Service and ServiceVersion can be automatically inferred when running
		// on App Engine.
		// ProjectID must be set if not running on GCP.
		// ProjectID: "my-project",
	}); err != nil {
		// TODO: Handle error.
	}
}

// [END profiler_setup_go_appengine]
