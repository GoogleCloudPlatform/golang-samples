// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package howto

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGetJob(t *testing.T) {
	buf := &bytes.Buffer{}
	if _, err := getJob(buf, testJob.Name); err != nil {
		t.Fatalf("getJob: %v", err)
	}
	want := testJob.Name
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("getJob got %q, want to contain %q", got, want)
	}
}

func TestUpdateJob(t *testing.T) {
	testJob.Qualifications = "foo"
	j, err := updateJob(ioutil.Discard, testJob.Name, testJob)
	if err != nil {
		t.Fatalf("updateJob: %v", err)
	}
	if j.Qualifications != "foo" {
		t.Fatalf("updateJob failed to update qualifications: got %q, want %q", j.Qualifications, "foo")
	}
	testJob.Qualifications = "bar"
	j, err = updateJob(ioutil.Discard, testJob.Name, testJob)
	if err != nil {
		t.Fatalf("updateJob: %v", err)
	}
	if j.Qualifications != "bar" {
		t.Fatalf("updateJob failed to update qualifications: got %q, want %q", j.Qualifications, "bar")
	}
}

func TestUpdateJobWithMask(t *testing.T) {
	testJob.Qualifications = "foo"
	j, err := updateJobWithMask(ioutil.Discard, testJob.Name, "qualifications", testJob)
	if err != nil {
		t.Fatalf("updateJobWithMask: %v", err)
	}
	if j.Qualifications != "foo" {
		t.Fatalf("updateJobWithMask failed to update qualifications: got %q, want %q", j.Qualifications, "foo")
	}
	testJob.Qualifications = "bar"
	j, err = updateJobWithMask(ioutil.Discard, testJob.Name, "qualifications", testJob)
	if err != nil {
		t.Fatalf("updateJobWithMask: %v", err)
	}
	if j.Qualifications != "bar" {
		t.Fatalf("updateJobWithMask failed to update qualifications: got %q, want %q", j.Qualifications, "bar")
	}
}

func TestListJobs(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := &bytes.Buffer{}
	if _, err := listJobs(buf, tc.ProjectID, fmt.Sprintf("companyName=%q", testCompany.Name)); err != nil {
		t.Fatalf("listJobs: %v", err)
	}
	want := testJob.Name
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("listJobs got %q, want to contain %q", got, want)
	}
}
