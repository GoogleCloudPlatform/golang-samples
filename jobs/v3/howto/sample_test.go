// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func checkServiceAvailable(t *testing.T, projectID string) {
	service, err := createCTSService()
	if err != nil {
		t.Skipf("createCTSService: service account likely in different project: %v", err)
	}
	if _, err := service.Projects.Companies.List("projects/" + projectID).Do(); err != nil {
		t.Skip("List: service account likely in different project")
	}
}

func TestRunBasicCompanySample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)
	out := new(bytes.Buffer)
	runBasicCompanySample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "GetCompany: Google Sample\n"
	want += "UpdateCompany: Google Sample (updated)\n"
	want += "UpdateCompanyWithMask: Google Sample (updated with mask)\n"
	want += "DeleteCompany StatusCode: 200\n"
	want += "ListCompanies Request ID:"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
}

func TestRunBasicJobSample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)

	out := new(bytes.Buffer)
	runBasicJobSample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Software Engineer\n"
	want += "GetJob: Software Engineer\n"
	want += "UpdateJob: Software Engineer (updated)\n"
	want += "UpdateJobWithMask: Title: Software Engineer (updated) Department: Engineering (updated with mask)\n"
	want += "ListJobs Request ID: "
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
	want = "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
}

func TestRunCommuteSearchSample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)
	out := new(bytes.Buffer)
	runCommuteSearchSample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Software Engineer\n"
	want += "CommuteSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job:"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
	want = "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
}

func TestRunHistogramSearchSample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)
	out := new(bytes.Buffer)
	runHistogramSearchSample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Software Engineer\n"
	want += "HistogramSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Software Engineer\n"
	want += "SimpleHistogramResults size: 1\n"
	want += "-- simple histogram searchType: COMPANY_ID value: "
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

	want = "CustomAttributeHistogramResults size: 1\n"
	want += "-- custom-attribute histogram key: someFieldString value: "
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

	want = "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

}

func TestRunEmailAlertSearchSample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)
	out := new(bytes.Buffer)
	runEmailAlertSearchSample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Software Engineer\n"

	want += "SearchForAlerts StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Software Engineer\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

	want = "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

}

func TestRunFeaturedJobSearchSample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)
	out := new(bytes.Buffer)
	runFeaturedJobSearchSample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Software Engineer (Featured)\n"

	want += "SearchFeaturedJobs StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Software Engineer (Featured)\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

	want = "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

}

func TestRunAutoCompleteSample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)
	out := new(bytes.Buffer)
	runAutoCompleteSample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Software Engineer\n"
	want += "CreateJob: GAP Product Manager\n"

	want += "DefaultAutoComplete query: sof StatusCode: 200\n"
	want += "-- suggestion: Software Engineer\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

	want = "JobTitleAutoComplete query: sof StatusCode: 200\n"
	want += "-- suggestion: Software Engineer\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

	want = "DefaultAutoComplete query: gap StatusCode: 200\n"
	want += "-- suggestion: Gap\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

	want = "JobTitleAutoComplete query: gap StatusCode: 200\n"
	//	want += "-- suggestion: GAP Product Manager\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

	want = "DeleteJob StatusCode: 200\n"
	want += "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

}

func TestRunCustomAttributeSample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)
	out := new(bytes.Buffer)
	runCustomAttributeSample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Software Engineer\n"

	want += "FilterOnStringValueCustomAttribute StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Software Engineer\n"
	want += "FilterOnLongValueCustomAttribute StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Software Engineer\n"
	want += "FilterOnMultiCustomAttributes StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Software Engineer\n"
	want += "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

}

func TestRunLocationBasedSearchSample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)
	out := new(bytes.Buffer)
	runLocationBasedSearchSample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Software Engineer\n"
	want += "CreateJob: Senior Software Engineer\n"

	want += "BasicLocationSearch StatusCode: 200\n"
	want += "MatchingJobs size: 2\n"
	want += "CityLocationSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Software Engineer\n"
	want += "BroadeningLocationSearch StatusCode: 200\n"
	want += "MatchingJobs size: 2\n"
	want += "KeywordLocationSearch StatusCode: 200\n"
	want += "MatchingJobs size: 2\n"
	want += "MultiLocationsSearch StatusCode: 200\n"
	want += "MatchingJobs size: 2\n"

	want += "DeleteJob StatusCode: 200\n"
	want += "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

}

func TestRunGeneralSearchSample(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)
	out := new(bytes.Buffer)
	runGeneralSearchSample(out, tc.ProjectID)
	got := out.String()

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Systems Administrator\n"

	want += "BasicJobSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Systems Administrator\n"
	want += "CategoryFilterSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Systems Administrator\n"
	want += "EmploymentTypesSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Systems Administrator\n"
	want += "DateRangeSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Systems Administrator\n"
	want += "LanguageCodeSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Systems Administrator\n"
	want += "CompanyDisplayNameSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Systems Administrator\n"
	want += "CompensationSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Systems Administrator\n"

	want += "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

}
