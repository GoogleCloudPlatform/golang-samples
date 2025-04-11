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
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGetCompany(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	buf := &bytes.Buffer{}
	if _, err := getCompany(buf, testCompany.Name); err != nil {
		t.Fatalf("getCompany: %v", err)
	}
	want := "Company: "
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("getCompany got %q, want %q", got, want)
	}
}

func TestUpdateCompany(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	testCompany.HiringAgency = false
	c, err := updateCompany(ioutil.Discard, testCompany.Name, testCompany)
	if err != nil {
		t.Fatalf("updateCompany: %v", err)
	}
	if c.HiringAgency {
		t.Fatalf("updateCompany did not set HiringAgency to false")
	}
	testCompany.HiringAgency = true
	c, err = updateCompany(ioutil.Discard, testCompany.Name, testCompany)
	if err != nil {
		t.Fatalf("updateCompany: %v", err)
	}
	if !c.HiringAgency {
		t.Fatalf("updateCompany did not set HiringAgency to true")
	}
}

func TestUpdateCompanyWithMask(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	testCompany.HiringAgency = false
	c, err := updateCompanyWithMask(ioutil.Discard, testCompany.Name, "hiring_agency", testCompany)
	if err != nil {
		t.Fatalf("updateCompanyWithMask: %v", err)
	}
	if c.HiringAgency {
		t.Fatalf("updateCompanyWithMask did not set HiringAgency to false")
	}
	testCompany.HiringAgency = true
	c, err = updateCompanyWithMask(ioutil.Discard, testCompany.Name, "hiring_agency", testCompany)
	if err != nil {
		t.Fatalf("updateCompanyWithMask: %v", err)
	}
	if !c.HiringAgency {
		t.Fatalf("updateCompanyWithMask did not set HiringAgency to true")
	}
}

func TestListCompanies(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	tc := testutil.SystemTest(t)

	buf := &bytes.Buffer{}
	if _, err := listCompanies(buf, tc.ProjectID); err != nil {
		t.Fatalf("listCompanies: %v", err)
	}
	want := testCompany.Name
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("listCompanies got %q, want to contain %q", got, want)
	}
}
