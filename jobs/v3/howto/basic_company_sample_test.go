// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package howto

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGetCompany(t *testing.T) {
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
