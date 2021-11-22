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

type EnvVars struct {
	// Job-defined env vars
	taskNum string
	attemptNum string

	// User-defined env vars
	sleepMs int64
	failRate float64
}

func stringToInt(s string) int64 {
	num, _ := strconv.ParseInt(s, 10, 64)
	return num
}

func stringToFloat(s string) float64 {
	num, _ := strconv.ParseFloat(s, 64)
	return num
}

func main() {
	env := &EnvVars{
		taskNum:     os.Getenv("TASK_NUM"),
		attemptNum:  os.Getenv("ATTEMPT_NUM"),
		sleepMs:   stringToInt(os.Getenv("SLEEP_MS")),
		failRate:  stringToFloat(os.Getenv("FAIL_RATE")),
	}

	log.Printf("Starting Task #%s, Attempt #%s ...", env.taskNum, env.attemptNum)

	// Simulate work
	if env.sleepMs > 0 {
		time.Sleep(time.Duration(env.sleepMs) * time.Millisecond)
	}

	// Simulate errors
	if env.failRate > 0 {
		if err := randomFailure(env); err != nil {
			log.Fatalf("%v", err)
		}
	}

	log.Printf("Completed Task #%s, Attempt #%s", env.taskNum, env.attemptNum)
}

// Throw an error based on fail rate
func randomFailure(env *EnvVars) error {
	if env.failRate < 0 || env.failRate > 1 {
		return fmt.Errorf("Invalid FAIL_RATE env var value: %f. Must be a float between 0 and 1 inclusive.", env.failRate)
	}

	rand.Seed(time.Now().UnixNano())
	randomFailure := rand.Float64()

	if randomFailure < env.failRate {
		return fmt.Errorf("Task #%s, Attempt #%s failed.", env.taskNum, env.attemptNum)
	}
	return nil
}
