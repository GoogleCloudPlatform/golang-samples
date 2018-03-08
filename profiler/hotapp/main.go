// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START hotapp]

// Sample hotapp is a synthetic application that exhibits different types of
// profiling hotspots: CPU, heap allocations, thread contention.
package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"cloud.google.com/go/profiler"
)

// There are several threads continuously fighting for this mutex.
var mu sync.Mutex

// Some allocated memory. Held in a global variable to protect it from GC.
var mem [][]byte

func sleepLocked(d time.Duration) {
	mu.Lock()
	time.Sleep(d)
	mu.Unlock()
}

func contentionImpl(d time.Duration) {
	for {
		mu.Lock()
		time.Sleep(d)
		mu.Unlock()
	}
}

func contention(d time.Duration) {
	contentionImpl(d)
}

func waitImpl() {
	select {}
}

// Waits forever.
func wait() {
	waitImpl()
}

func allocImpl() {
	// Allocate 64 MiB in 64 KiB chunks
	for i := 0; i < 64*16; i++ {
		mem = append(mem, make([]byte, 64*1024))
	}
}

func alloc() {
	allocImpl()
}

func busyloop() {
	for {
		load()
		// Yield so that some preemption happens.
		runtime.Gosched()
	}
}

func load() {
	for i := 0; i < (1 << 20); i++ {
	}
}

func main() {
	err := profiler.Start(profiler.Config{
		Service:        "hotapp",
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
	alloc()

	// Simulate CPU load.
	busyloop()
}

// [END hotapp]
