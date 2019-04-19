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
