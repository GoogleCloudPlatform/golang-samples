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
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}
	env := map[string]string{}

	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		stdOut, stdErr, err := m.Run(env, 2*time.Minute, "-project", projectID)

		if err != nil {
			r.Errorf("execution failed: %v, stdOut:%s, stdErr:%s", err, stdOut, stdErr)
			return
		}
		if len(stdErr) > 0 {
			r.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
		}
		got := string(stdOut)
		want := "map[born:1912 first:Alan last:Turing middle:Mathison]"
		if !strings.Contains(got, want) {
			r.Errorf("stdout returned %s, wanted to contain %s", got, want)
		}
	})
}
