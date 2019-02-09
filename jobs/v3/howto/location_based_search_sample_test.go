// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package howto

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestBasicLocationSearch(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := basicLocationSearch(buf, tc.ProjectID, testCompany.Name, "Mountain View, CA", .5); err != nil {
			r.Errorf("basicLocationSearch: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("basicLocationSearch got %q, want to contain %q", got, want)
		}
	})
}

func TestCityLocationSearch(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := cityLocationSearch(buf, tc.ProjectID, testCompany.Name, "Mountain View, CA"); err != nil {
			r.Errorf("cityLocationSearch: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("cityLocationSearch got %q, want to contain %q", got, want)
		}
	})
}

func TestBroadeningLocationSearch(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := broadeningLocationSearch(buf, tc.ProjectID, testCompany.Name, "Sunnyvale, CA"); err != nil {
			r.Errorf("broadeningLocationSearch: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("broadeningLocationSearch got %q, want to contain %q", got, want)
		}
	})
}

func TestKeywordLocationSearch(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := keywordLocationSearch(buf, tc.ProjectID, testCompany.Name, "Mountain View, CA", .5, "SWE"); err != nil {
			r.Errorf("keywordLocationSearch: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("keywordLocationSearch got %q, want to contain %q", got, want)
		}
	})
}

func TestMultiLocationsSearch(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := multiLocationsSearch(buf, tc.ProjectID, testCompany.Name, "New York, NY", "Sunnyvale, CA", .5); err != nil {
			r.Errorf("multiLocationsSearch: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("multiLocationsSearch got %q, want to contain %q", got, want)
		}
	})
}
