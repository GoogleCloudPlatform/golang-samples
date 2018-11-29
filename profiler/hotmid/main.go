// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START hotmid]

// Sample hotmid is an application that simulates multiple calls to a library
// function made via different call paths. Each of these calls is not
// particularly expensive (and so does not stand out on the flame graph). But
// in the aggregate these calls add up to a significant time which can be
// identified via looking at the flat list of functions' self and total time.
package main

import (
	"flag"
	"log"
	"time"

	"cloud.google.com/go/profiler"
)

var (
	// version is the service version to configure.
	version = flag.String("version", "1.0.0", "service version")
	// seconds is the benchmark duration in seconds or 0 to run forever.
	seconds = flag.Int("seconds", 0, "benchmark duration in seconds")
)

const spinCount = 1e6 // Ballpark to spin for ~0.5ms.

// Simulates a logging function that is called all over the place by the
// middleware. It illustrates the case of a function that is called by multiple
// call paths so it isn't immediately obvious as a hotspot when looking at a
// flame graph. Also note that the function itself does not have significant
// self time since it just calls out to other functions. So, to identify the
// function at its true cost, aggregation by total time is needed.
func myLog() {
	myLogMutexLock()
	myLogWrite()
	myLogBuffer()
}

func myLogMutexLock() {
	for i := 0; i < spinCount; i++ {
	}
}

func myLogBuffer() {
	for i := 0; i < 2*spinCount; i++ {
	}
}

func myLogWrite() {
	for i := 0; i < 5*spinCount; i++ {
	}
}

func processItem(depth int) {
	if depth == 0 {
		return
	}
	for i := 0; i < 10*spinCount; i++ {
	}
	myLog()
	processItem(depth - 1)
}

func foo1() {
	processItem(5)
}

func foo2() {
	processItem(5)
}

func foo3() {
	processItem(10)
}

func foo4() {
	processItem(10)
}

func foo5() {
	processItem(15)
}

func foo6() {
	processItem(15)
}

func run() {
	start, duration := time.Now(), time.Duration(*seconds)*time.Second
	for duration == 0 || time.Since(start) < duration {
		foo1()
		foo2()
		foo3()
		foo4()
		foo5()
		foo6()
	}
}

func main() {
	flag.Parse()

	err := profiler.Start(profiler.Config{
		Service:        "hotmid-service",
		ServiceVersion: *version,
		DebugLogging:   true,
	})
	if err != nil {
		log.Fatalf("failed to start the profiler: %v", err)
	}

	run()
}

// [END hotmid]
