// Copyright 2023 Google LLC
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

package jobs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestListJobs(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	listJobs(buf, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
	s := buf.String()
	if len(s) == 0 {
		// Create job.
		riskNumerical(tc.ProjectID, "bigquery-public-data", "risk-topic", "risk-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		buf.Reset()
		err := listJobs(buf, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
		if err != nil {
			t.Errorf("listJobs(%s, %s, %s) = error %q, want nil", buf, tc.ProjectID, "", err)
		}
		s = buf.String()
	}
	if !strings.Contains(buf.String(), "Job") {
		t.Errorf("%q not found in listJobs output: %q", "Job", s)
	}
}
