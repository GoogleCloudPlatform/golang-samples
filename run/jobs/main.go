// Copyright 2021 Google LLC
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
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Retrieve Job-defined env vars
var TASK_NUM = os.Getenv("TASK_NUM")
var ATTEMPT_NUM = os.Getenv("ATTEMPT_NUM")

// Retrieve User-defined env vars
var SLEEP_MS = os.Getenv("SLEEP_MS")
var FAIL_RATE = os.Getenv("FAIL_RATE")

func main() {
	log.Printf("Starting Task #%s, Attempt #%s ...", TASK_NUM, ATTEMPT_NUM)

	// Simulate work
	if SLEEP_MS != "" {
		// Convert SLEEP_MS from String to Int
		SLEEP_MS, _ := strconv.Atoi(SLEEP_MS)
		time.Sleep(time.Duration(SLEEP_MS) * time.Millisecond)
	}

	// Simulate errors
	if FAIL_RATE != "" {
		if err := randomFailure(FAIL_RATE); err != nil {
			log.Fatalf("%v", err)
		}
	}

	log.Printf("Completed Task #%s, Attempt #%s", TASK_NUM, ATTEMPT_NUM)
}

// Throw an error based on fail rate
func randomFailure(FAIL_RATE string) error {
	rate, err := strconv.ParseFloat(FAIL_RATE, 64)

	if err != nil || rate < 0 || rate > 1 {
		return fmt.Errorf("Invalid FAIL_RATE env var value: %s. Must be a float between 0 and 1 inclusive.", FAIL_RATE)
	}

	rand.Seed(time.Now().UnixNano())
	randomFailure := rand.Float64()

	if randomFailure < rate {
		return fmt.Errorf("Task #%s, Attempt #%s failed.", TASK_NUM, ATTEMPT_NUM)
	}
	return nil
}
