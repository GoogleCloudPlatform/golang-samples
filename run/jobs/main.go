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

// [START cloudrun_jobs_quickstart]
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Job-defined
	taskNum    string
	attemptNum string

	// User-defined
	sleepMs  int64
	failRate float64
}

func configFromEnv() (Config, error) {
	// [START cloudrun_jobs_env_vars]
	// Job-defined
	taskNum := os.Getenv("CLOUD_RUN_TASK_INDEX")
	attemptNum := os.Getenv("CLOUD_RUN_TASK_ATTEMPT")
	// User-defined
	sleepMs, err := sleepMsToInt(os.Getenv("SLEEP_MS"))
	failRate, err := failRateToFloat(os.Getenv("FAIL_RATE"))
	// [END cloudrun_jobs_env_vars]

	if err != nil {
		return Config{}, err
	}

	config := Config{
		taskNum:    taskNum,
		attemptNum: attemptNum,
		sleepMs:    sleepMs,
		failRate:   failRate,
	}
	return config, nil
}

func sleepMsToInt(s string) (int64, error) {
	sleepMs, err := strconv.ParseInt(s, 10, 64)
	return sleepMs, err
}

func failRateToFloat(s string) (float64, error) {
	// Default empty variable to 0
	if s == "" {
		return 0, nil
	}

	// Convert string to float
	failRate, err := strconv.ParseFloat(s, 64)

	// Check that rate is valid
	if failRate < 0 || failRate > 1 {
		return failRate, fmt.Errorf("Invalid FAIL_RATE value: %f. Must be a float between 0 and 1 inclusive.", failRate)
	}

	return failRate, err
}

func main() {
	config, err := configFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting Task #%s, Attempt #%s ...", config.taskNum, config.attemptNum)

	// Simulate work
	if config.sleepMs > 0 {
		time.Sleep(time.Duration(config.sleepMs) * time.Millisecond)
	}

	// Simulate errors
	if config.failRate > 0 {
		if failure := randomFailure(config); failure != nil {
			// [START cloudrun_jobs_exit_process]
			log.Fatalf("%v", failure)
			// [END cloudrun_jobs_exit_process]
		}
	}

	log.Printf("Completed Task #%s, Attempt #%s", config.taskNum, config.attemptNum)
}

// Throw an error based on fail rate
func randomFailure(config Config) error {
	rand.Seed(time.Now().UnixNano())
	randomFailure := rand.Float64()

	if randomFailure < config.failRate {
		return fmt.Errorf("Task #%s, Attempt #%s failed.", config.taskNum, config.attemptNum)
	}
	return nil
}

// [END cloudrun_jobs_quickstart]
