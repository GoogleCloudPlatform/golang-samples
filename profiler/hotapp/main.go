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

// [START hotapp]

// Sample hotapp is a synthetic application that exhibits different types of
// profiling hotspots: CPU, heap allocations, thread contention.
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
	version = flag.String("version", "1.0.0", "service version")
	// Skew of foo1 function over foo2, in the CPU busyloop, to simulate diff.
	skew = flag.Int("skew", 100, "skew of foo2 over foo1: foo2 will consume skew/100 CPU time compared to foo1 (default is no skew)")
	// There are several goroutines continuously fighting for this mutex.
	mu sync.Mutex
	// Some allocated memory. Held in a global variable to protect it from GC.
	mem [][]byte
)

func sleepLocked(d time.Duration) {
	mu.Lock()
	time.Sleep(d)
	mu.Unlock()
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
		foo1(100)
		foo2(*skew)
		// Yield so that some preemption happens.
		runtime.Gosched()
	}
}

func foo1(scale int) {
	bar(scale)
	baz(scale)
}

func foo2(scale int) {
	bar(scale)
	baz(scale)
}

func bar(scale int) {
	load(scale)
}

func baz(scale int) {
	load(scale)
}

func load(scale int) {
	for i := 0; i < scale*(1<<16); i++ {
	}
}

func main() {
	flag.Parse()

	err := profiler.Start(profiler.Config{
		Service:        "hotapp-service",
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

// [END hotapp]
