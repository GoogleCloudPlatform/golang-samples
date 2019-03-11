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

package howto

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGetJob(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := strings.SplitAfter(testJob.Name, "jobs/")[1]
	if _, err := getJob(buf, tc.ProjectID, jobID); err != nil {
		t.Fatalf("getJob: %v", err)
	}
	want := "Software Engineer"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("getJob got %q, want to contain %q", got, want)
	}
}
