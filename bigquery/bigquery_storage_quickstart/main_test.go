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

package main

import (
	"fmt"
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
		// TODO(shollyman): remove after https://github.com/GoogleCloudPlatform/golang-samples/issues/1866 deflaked.
		// We're running over the 30s timeout, but unclear why; normal execution is sub-5s.
		t.Logf("stderr: %s", string(stdErr))
	}

	// We don't look for specific strings, just expect at least 1kb of output.
	if len(stdOut) < 1024 {
		t.Errorf("expected more output.  Stdout: %s", string(stdOut))
	}

	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
}
