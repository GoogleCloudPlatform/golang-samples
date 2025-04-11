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
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/talent/apiv4beta1/talentpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
)

var testCompany *talentpb.Company
var testJob *talentpb.Job

func TestBasicUsage(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	tc := testutil.SystemTest(t)

	externalID := fmt.Sprintf("company-%s", uuid.Must(uuid.NewV4()).String())
	displayName := "Google Sample"
	var err error

	// Create the company.
	testCompany, err = createCompany(ioutil.Discard, tc.ProjectID, externalID, displayName)
	if err != nil {
		log.Fatalf("createCompany: %v", err)
	}

	companyID := strings.SplitAfter(testCompany.Name, "companies/")[1]
	requisitionID := fmt.Sprintf("job-%s", uuid.Must(uuid.NewV4()).String())
	title := "Software Engineer"
	URI := "https://googlesample.com/career"
	description := "Design, develop, test, deploy, maintain and improve software."
	address1 := "1600 Amphitheatre Pkwy, Mountain View, CA 94043"
	address2 := "85 10th Ave, New York, NY 10011"
	languageCode := "en-US"

	// Create the sample job, exercised by various search examples.
	testJob, err = createJob(ioutil.Discard, tc.ProjectID, companyID, requisitionID, title, URI, description, address1, address2, languageCode)
	if err != nil {
		log.Fatalf("createJob: %v", err)
	}
	jobID := strings.SplitAfter(testJob.Name, "jobs/")[1]
	buf := &bytes.Buffer{}

	// Test company listing.
	listCompanies(buf, tc.ProjectID)
	if err := listCompanies(buf, tc.ProjectID); err != nil {
		t.Fatalf("listCompanies: %v", err)
	}
	want := testCompany.Name
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("listCompanies got %q, want to contain %q", got, want)
	}

	// List jobs.
	buf.Reset()
	if err := listJobs(buf, tc.ProjectID, companyID); err != nil {
		t.Fatalf("listJobs: %v", err)
	}
	want = testJob.Name
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("listJobs got %q, want to contain %q", got, want)
	}

	// Get a single company.
	buf.Reset()
	if _, err := getCompany(buf, tc.ProjectID, companyID); err != nil {
		t.Fatalf("getCompany: %v", err)
	}
	want = "Google"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("getCompany got %q, want %q", got, want)
	}

	// Get a single job.
	buf.Reset()
	if _, err := getCompany(buf, tc.ProjectID, companyID); err != nil {
		t.Fatalf("getCompany: %v", err)
	}
	want = "Google"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("getCompany got %q, want %q", got, want)
	}

	// Client event testing.
	buf.Reset()
	requestID := fmt.Sprintf("requestId-%s", uuid.Must(uuid.NewV4()).String())
	eventID := fmt.Sprintf("eventId-%s", uuid.Must(uuid.NewV4()).String())
	relatedJobNames := []string{testJob.Name}

	if _, err := createClientEvent(buf, tc.ProjectID, requestID, eventID, relatedJobNames); err != nil {
		log.Fatalf("createClientEvent: %v", err)
	}
	want = "Client event created: "
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("getJob got %q, want to contain %q", got, want)
	}

	// Exercise various search tests.
	buf.Reset()
	if _, err := jobTitleAutocomplete(buf, tc.ProjectID, "Software"); err != nil {
		t.Fatalf("jobTitleAutoComplete: %v", err)
	}
	want = "Software Developer"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("jobTitleAutoComplete got %q, want %q", got, want)
	}

	// Commute search is problematic, wrapped in a retry to reduce flakiness.
	testutil.Retry(t, 20, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := commuteSearch(buf, tc.ProjectID, companyID); err != nil {
			r.Errorf("commuteSearch: %v", err)
			return
		}
		want := "Mountain View"
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("commuteSearch got %q, want to contain %q", got, want)
		}
	})

	// Histogram search.
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := histogramSearch(buf, tc.ProjectID, companyID); err != nil {
			r.Errorf("histogramSearch: %v", err)
		}
		want := strings.SplitAfter(testJob.Name, "jobs/")[1]
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("histogramSearch got %q, want to contain %q", got, want)
		}
	})

	// Cleanup job and company.
	if err := deleteJob(ioutil.Discard, tc.ProjectID, jobID); err != nil {
		log.Fatalf("deleteJob: %v", err)
	}

	if err := deleteCompany(ioutil.Discard, tc.ProjectID, companyID); err != nil {
		log.Fatalf("deleteCompany: %v", err)
	}
}

func TestBatchDeleteJobs(t *testing.T) {
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		tc := testutil.SystemTest(t)

		requisitionID := fmt.Sprintf("job-%s", uuid.Must(uuid.NewV4()).String())
		title := "Software Engineer"
		URI := "https://googlesample.com/career"
		description := "Design, develop, test, deploy, maintain and improve software."
		address1 := "Mountain View, CA"
		address2 := "New York City, NY"
		languageCode1 := "en-US"
		languageCode2 := "sr_Latn"

		// Company setup.
		externalID := fmt.Sprintf("company-%s", uuid.Must(uuid.NewV4()).String())
		displayName := "Google Sample"
		testCompany, err := createCompany(ioutil.Discard, tc.ProjectID, externalID, displayName)
		if err != nil {
			log.Fatalf("createCompany: %v", err)
		}
		companyID := strings.SplitAfter(testCompany.Name, "companies/")[1]

		// Create two identical jobs with different language codes.
		if _, err := createJob(ioutil.Discard, tc.ProjectID, companyID, requisitionID, title, URI, description, address1, address2, languageCode1); err != nil {
			log.Fatalf("createJob1: %v", err)
		}
		if _, err := createJob(ioutil.Discard, tc.ProjectID, companyID, requisitionID, title, URI, description, address1, address2, languageCode2); err != nil {
			log.Fatalf("createJob2: %v", err)
		}

		if err := batchDeleteJobs(ioutil.Discard, tc.ProjectID, companyID, requisitionID); err != nil {
			log.Fatalf("batchDeleteJob: %v", err)
		}

		if err := deleteCompany(ioutil.Discard, tc.ProjectID, companyID); err != nil {
			log.Fatalf("deleteCompany: %v", err)
		}
	})
}

func TestCreateJobWithCustomAttributes(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")
	tc := testutil.SystemTest(t)
	// Company setup.
	externalID := fmt.Sprintf("company-%s", uuid.Must(uuid.NewV4()).String())
	displayName := "Google Sample"
	testCompany, err := createCompany(ioutil.Discard, tc.ProjectID, externalID, displayName)
	if err != nil {
		log.Fatalf("createCompany: %v", err)
	}
	companyID := strings.SplitAfter(testCompany.Name, "companies/")[1]
	buf := &bytes.Buffer{}

	// Create job.
	customJob, err := createJobWithCustomAttributes(buf, tc.ProjectID, companyID, testJob.Title)
	if err != nil {
		log.Fatalf("createJobWithCustomAttributes: %v", err)
	}
	want := "900"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("createJobWithCustomAttributes got %q, want to contain %q", got, want)
	}

	testutil.Retry(t, 20, 2*time.Second, func(r *testutil.R) {
		// Custom ranking search.
		buf.Reset()
		if err := customRankingSearch(buf, tc.ProjectID, companyID); err != nil {
			r.Errorf("customRankingSearch: %v", err)
		}
		want = "Job: "
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("getJob got %q, want to contain %q", got, want)
		}

	})

	jobID := strings.SplitAfter(customJob.Name, "jobs/")[1]
	if err := deleteJob(ioutil.Discard, tc.ProjectID, jobID); err != nil {
		log.Fatalf("deleteJob: %v", err)
	}
}
