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
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestSuccessfulJob(t *testing.T) {
	os.Setenv("SLEEP_MS", "0")
	os.Unsetenv("FAIL_RATE")
	os.Setenv("CLOUD_RUN_TASK_INDEX", "1")
	os.Setenv("CLOUD_RUN_TASK_ATTEMPT", "1")

	var buf bytes.Buffer
	log.SetOutput(&buf)

	main()
	log.SetOutput(os.Stdout)
	output := buf.String()

	start := "Starting Task #1, Attempt #1 ..."
	finish := "Completed Task #1, Attempt #1"

	if !(strings.Contains(output, start)) {
		t.Errorf("\nWant:\n%s\n\nGot:\n%s", start, output)
	}

	if !(strings.Contains(output, finish)) {
		t.Errorf("\nWant:\n%s\n\nGot:\n%s", finish, output)
	}
}

func TestRandomFailure(t *testing.T) {
	config := Config{
		taskNum:    "1",
		attemptNum: "1",
		sleepMs:    2,
		failRate:   1,
	}

	err := randomFailure(config)
	if err == nil {
		t.Fatalf("Test should fail with FAIL_RATE 1")
	}

	config.failRate = 0
	err = randomFailure(config)
	if err != nil {
		t.Fatalf("Test should pass with empty FAIL_RATE")
	}
}
