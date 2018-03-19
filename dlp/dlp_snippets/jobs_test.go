// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
)

func TestListJobs(t *testing.T) {
	buf := new(bytes.Buffer)
	listJobs(buf, client, projectID, "", "RISK_ANALYSIS_JOB")
	s := buf.String()
	if len(s) == 0 {
		// Create job.
		riskNumerical(ioutil.Discard, client, projectID, "bigquery-public-data", "risk-topic", "risk-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		buf.Reset()
		listJobs(buf, client, projectID, "", "RISK_ANALYSIS_JOB")
		s = buf.String()
	}
	if !strings.Contains(buf.String(), "Job") {
		t.Errorf("%q not found in listJobs output: %q", "Job", s)
	}
}

var jobIDRegexp = regexp.MustCompile(`Job ([^ ]+) status.*`)

func TestDeleteJob(t *testing.T) {
	buf := new(bytes.Buffer)
	listJobs(buf, client, projectID, "", "RISK_ANALYSIS_JOB")
	s := buf.String()
	if len(s) == 0 {
		// Create job.
		riskNumerical(ioutil.Discard, client, projectID, "bigquery-public-data", "risk-topic", "risk-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		buf.Reset()
		listJobs(buf, client, projectID, "", "RISK_ANALYSIS_JOB")
		s = buf.String()
	}
	jobName := string(jobIDRegexp.FindSubmatch([]byte(s))[1])
	buf.Reset()
	deleteJob(buf, client, jobName)
	if got := buf.String(); got != "Successfully deleted job" {
		t.Errorf("unable to delete job: %s", s)
	}
}
