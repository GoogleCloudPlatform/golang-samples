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
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateJobWithCustomAttributes(t *testing.T) {
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		tc := testutil.SystemTest(t)
		buf := &bytes.Buffer{}
		companyID := strings.SplitAfter(testCompany.Name, "companies/")[1]
		customJob, err := createJobWithCustomAttributes(buf, tc.ProjectID, companyID, testJob.Title)
		if err != nil {
			log.Fatalf("createJobWithCustomAttributes: %v", err)
		}
		want := "900"
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("createJobWithCustomAttributes got %q, want to contain %q", got, want)
		}
		jobID := strings.SplitAfter(customJob.Name, "jobs/")[1]
		if err := deleteJob(ioutil.Discard, tc.ProjectID, jobID); err != nil {
			log.Fatalf("deleteJob: %v", err)
		}
	})
}
