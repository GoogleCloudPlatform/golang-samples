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

// Package jobs contains example snippets using the DLP jobs API.
package jobs

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	dlp "cloud.google.com/go/dlp/apiv2"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var client *dlp.Client
var projectID string

func TestMain(m *testing.M) {
	ctx := context.Background()
	if c, ok := testutil.ContextMain(m); ok {
		var err error
		client, err = dlp.NewClient(ctx)
		if err != nil {
			log.Fatalf("datastore.NewClient: %v", err)
		}
		projectID = c.ProjectID
		defer client.Close()
	}
	os.Exit(m.Run())
}

func TestListJobs(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	listJobs(buf, client, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
	s := buf.String()
	if len(s) == 0 {
		// Create job.
		riskNumerical(ioutil.Discard, client, tc.ProjectID, "bigquery-public-data", "risk-topic", "risk-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		buf.Reset()
		err := listJobs(buf, client, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
		if err != nil {
			t.Errorf("listJobs(%s, %s, %s) = error %q, want nil", buf, client, tc.ProjectID, err)
		}
		s = buf.String()
	}
	if !strings.Contains(buf.String(), "Job") {
		t.Errorf("%q not found in listJobs output: %q", "Job", s)
	}
}

var jobIDRegexp = regexp.MustCompile(`Job ([^ ]+) status.*`)

func TestDeleteJob(t *testing.T) {
	testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	listJobs(buf, client, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
	s := buf.String()
	if len(s) == 0 {
		// Create job.
		riskNumerical(ioutil.Discard, client, tc.ProjectID, "bigquery-public-data", "risk-topic", "risk-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		buf.Reset()
		listJobs(buf, client, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
		s = buf.String()
	}
	jobName := string(jobIDRegexp.FindSubmatch([]byte(s))[1])
	buf.Reset()
	deleteJob(buf, client, jobName)
	if got := buf.String(); got != "Successfully deleted job" {
		t.Errorf("unable to delete job: %s", s)
	}
}
