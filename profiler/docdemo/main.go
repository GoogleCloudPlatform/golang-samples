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

// [START profiler_docdemo]

// Sample docdemo is a synthetic application that exhibits different types of
// profiling hotspots: CPU, heap allocations, thread contention.
// Sample docemo is a modification of the sample hotapp adds work, in the
// form of for-loops, to the following methods: foo1, foo2, bar, baz.
package main

import (
	"flag"
	"log"
	"runtime"
	"sync"
	"time"

	"cloud.google.com/go/profiler"
)

var (
	// Service version to configure.
	version = flag.String("version", "1.75.0", "service version")
	// Skew changes the relative CPU time of foo1 to foo2.
	// Decreasing skew increases the ratio of the CPU time of foo1 to that of foo2.
	skew = flag.Int("skew", 75, "Set skew from 0 to 100. Decreasing skew reduces relative CPU time for foo2.")
	// There are several goroutines continuously fighting for this mutex.
	mu sync.Mutex
	// Some allocated memory. Held in a global variable to protect it from GC.
	mem [][]byte
)

func main() {
	flag.Parse()

	err := profiler.Start(profiler.Config{
		Service:        "docdemo-service",
		ServiceVersion: *version,
		DebugLogging:   true,
		MutexProfiling: true,
	})
	if err != nil {
		log.Fatalf("failed to start the profiler: %v", err)
	}

	// Use four OS threads for the contention simulation.
	runtime.GOMAXPROCS(4)
	for i := 0; i < 4; i++ {
		go contention(time.Duration(i) * 50 * time.Millisecond)
	}

	// Simulate some waiting goroutines.
	for i := 0; i < 100; i++ {
		go wait()
	}

	// Simulate some memory allocation.
	allocOnce()

	// Simulate repeated memory allocations.
	go allocMany()

	// Simulate CPU load.
	busyloop()
}

// Simulates some work that contends over a shared mutex. It calls an "impl"
// function to produce a bit deeper stacks in the profiler visualization,
// merely for illustration purpose.
func contention(d time.Duration) {
	contentionImpl(d)
}

func contentionImpl(d time.Duration) {
	for {
		mu.Lock()
		time.Sleep(d)
		mu.Unlock()
	}
}

// Waits forever simulating a goroutine that is consistently blocked on I/O.
// It calls an "impl" function to produce a bit deeper stacks in the profiler
// visualization, merely for illustration purpose.
func wait() {
	waitImpl()
}

func waitImpl() {
	select {}
}

// Simulates a memory-hungry function. It calls an "impl" function to produce
// a bit deeper stacks in the profiler visualization, merely for illustration
// purpose.
func allocOnce() {
	allocImpl()
}

func allocImpl() {
	// Allocate 64 MiB in 64 KiB chunks
	for i := 0; i < 64*16; i++ {
		mem = append(mem, make([]byte, 64*1024))
	}
}

// allocMany simulates a function which allocates a lot of memory, but does not
// hold on to that memory.
func allocMany() {
	// Allocate 1 MiB of 64 KiB chunks repeatedly.
	for {
		for i := 0; i < 16; i++ {
			_ = make([]byte, 64*1024)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Simulates a CPU-intensive computation.
func busyloop() {
	for {
		for i := 0; i < 100*(1<<16); i++ {
		}
		foo1(100)
		foo2(*skew)
		// Yield so that some preemption happens.
		runtime.Gosched()
	}
}

func foo1(scale int) {
	// local work
	for i := 0; i < scale*(1<<16); i++ {
	}
	bar(scale)
	baz(scale)
}

func foo2(scale int) {
	// local work
	for i := 0; i < 5*scale*(1<<16); i++ {
	}
	bar(scale)
	baz(scale)
}

func bar(scale int) {
	// local work
	for i := 0; i < scale*(1<<16); i++ {
	}
	load(scale)
}

func baz(scale int) {
	// local work
	for i := 0; i < 5*scale*(1<<16); i++ {
	}
	load(scale)
}

func load(scale int) {
	for i := 0; i < scale*(1<<16); i++ {
	}
}

// [END profiler_docdemo]
