// Copyright 2024 Google LLC
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

package snippets

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateJobWithCustomJobLabels(t *testing.T) {
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	jobName := fmt.Sprintf("test-job-go-%v-%v", time.Now().Format("2006-01-02"), r.Int())
	region := "us-central1"

	buf := &bytes.Buffer{}

	job, err := createJobWithCustomJobLabels(buf, tc.ProjectID, region, jobName)

	if err != nil {
		t.Errorf("createJobWithCustomJobLabels got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Job created") {
		t.Errorf("createJobWithCustomJobLabels got %q, expected %q", got, "Job created")
	}

	labels := job.GetLabels()
	if labels["env"] != "dev" || labels["type"] != "single command" {
		t.Errorf("labels weren't set")
	}

	if err := deleteJob(buf, tc.ProjectID, region, jobName); err != nil {
		t.Errorf("deleteJob got err: %v", err)
	}
}
