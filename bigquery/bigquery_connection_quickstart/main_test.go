// Copyright 2020 Google LLC
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
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestApp(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}

	stdOut, stdErr, err := m.Run(nil, 30*time.Second, fmt.Sprintf("--project_id=%s", tc.ProjectID))
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}

	// Look for a known substring in the output
	if !strings.Contains(string(stdOut), " connections processed.") {
		t.Errorf("Did not find expected output.  Stdout: %s", string(stdOut))
	}

	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
}
