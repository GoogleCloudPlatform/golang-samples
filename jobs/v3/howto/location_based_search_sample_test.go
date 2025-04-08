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
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestBasicLocationSearch(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

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
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

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
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := broadeningLocationSearch(buf, tc.ProjectID, testCompany.Name, "Bay Area"); err != nil {
			r.Errorf("broadeningLocationSearch: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("broadeningLocationSearch got %q, want to contain %q", got, want)
		}
	})
}

func TestKeywordLocationSearch(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

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
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

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
