// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START profiler_quickstart]

// Sample profiler_quickstart simulates a CPU-intensive workload for profiler.
package main

import (
	"log"
	"runtime"

	"cloud.google.com/go/profiler"
)

func busyloop() {
	for {
		load()
		// Make sure to yield so that the profiler thread
		// gets some CPU time even on single core machines
		// where GOMAXPROCS is 1. Not needed in real programs
		// as typically the preemption happens naturally.
		runtime.Gosched()
	}
}

func load() {
	for i := 0; i < (1 << 20); i++ {
	}
}

func main() {
	err := profiler.Start(profiler.Config{
		Service:              "hello-profiler",
		NoHeapProfiling:      true,
		NoAllocProfiling:     true,
		NoGoroutineProfiling: true,
		DebugLogging:         true,
		// ProjectID must be set if not running on GCP.
		// ProjectID: "my-project",
	})
	if err != nil {
		log.Fatalf("failed to start the profiler: %v", err)
	}
	busyloop()
}

// [END profiler_quickstart]
