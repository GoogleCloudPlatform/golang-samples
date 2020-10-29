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
	// "fmt"
	// "log"
	// "strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	// testResourceName := "my-resource-name"
	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	// testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
	_, stdErr, err := m.Run(nil, 2*time.Minute,
		"--project_id", tc.ProjectID,
		// "--resource_name", testResourceName,
	)
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}
	if len(stdErr) > 0 {
		// r.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
	// })
}
