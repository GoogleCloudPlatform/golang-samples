// Copyright 2022 Google LLC
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

package cloudruntests

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCloudRunJobs(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	crj := &cloudrunci.Job{
		Name:        "runjobs",
		ProjectID:   tc.ProjectID,
		Dir:         "../jobs",
		AsBuildpack: true,
		Region:      "us-central1",
		Env: map[string]string{
			"FAIL_RATE": "0.0",
			"SLEEP_MS":  "10000",
		},
	}

	if err := crj.Create(); err != nil {
		t.Fatalf("Create %q: %v", crj.Name, err)
	}
	if err := crj.Run(); err != nil {
		t.Errorf("Run(%s): %s", crj.Name, err)
	}

	found, err := crj.LogEntries("", "Completed Task", 5)
	if err != nil {
		t.Errorf("LogEntries: %v", err)
	}
	if !found {
		t.Errorf("Failed to find log entries for job")
	}

	defer crj.Clean()

}
