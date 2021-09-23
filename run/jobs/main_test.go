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
	"testing"
)

func ExampleSuccessfulRun() {
	os.Setenv("SLEEP_MS", "0")
	os.Setenv("TASK_NUM", "1")
	os.Setenv("ATTEMPT_NUM", "1")
	os.Setenv("FAIL_RATE", "")

	var buf bytes.Buffer
	log.SetOutput(&buf)

	main()
	log.SetOutput(os.Stdout)
	log.Print(buf.String())

	// Output:
	// Started Task #1, Attempt #1
	// Completed Task #1, Attempt #1
}

func TestRandomFailure(t *testing.T) {
	err := randomFailure("1")
	if err == nil {
		t.Fatalf("Test should fail with FAIL_RATE 1")
	}

	err = randomFailure("0")
	if err != nil {
		t.Fatalf("Test should pass with empty FAIL_RATE")
	}
}
