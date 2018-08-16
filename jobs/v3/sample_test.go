package cjdsample

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestBasicCompanySampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	BasicCompanySampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

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

func TestBasicJobSampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	BasicJobSampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

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

func TestCommuteSearchSampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	CommuteSearchSampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

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

func TestHistogramSearchSampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	HistogramSearchSampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

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

func TestEmailAlertSearchSampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	EmailAlertSearchSampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

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

func TestFeaturedJobSearchSampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	FeaturedJobSearchSampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

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

func TestAutoCompleteSampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	AutoCompleteSampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

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
	want += "-- suggestion: Gap Inc.\n"
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

func TestCustomAttributeSampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	CustomAttributeSampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

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

func TestLocationBasedSearchSampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	LocationBasedSearchSampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

	want := "CreateCompany: Google Sample\n"
	want += "CreateJob: Software Engineer\n"
	want += "CreateJob: Senior Software Engineer\n"

	want += "BasicLocationSearch StatusCode: 200\n"
	want += "MatchingJobs size: 2\n"
	want += "-- match job: Software Engineer\n"
	want += "-- match job: Senior Software Engineer\n"
	want += "CityLocationSearch StatusCode: 200\n"
	want += "MatchingJobs size: 1\n"
	want += "-- match job: Software Engineer\n"
	want += "BroadeningLocationSearch StatusCode: 200\n"
	want += "MatchingJobs size: 2\n"
	want += "-- match job: Software Engineer\n"
	want += "-- match job: Senior Software Engineer\n"
	want += "KeywordLocationSearch StatusCode: 200\n"
	want += "MatchingJobs size: 2\n"
	want += "-- match job: Software Engineer\n"
	want += "-- match job: Senior Software Engineer\n"
	want += "MultiLocationsSearch StatusCode: 200\n"
	want += "MatchingJobs size: 2\n"
	want += "-- match job: Software Engineer\n"
	want += "-- match job: Senior Software Engineer\n"

	want += "DeleteJob StatusCode: 200\n"
	want += "DeleteJob StatusCode: 200\n"
	want += "DeleteCompany StatusCode: 200\n"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}

}

func TestGeneralSearchSampleEntry(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	GeneralSearchSampleEntry()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

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
