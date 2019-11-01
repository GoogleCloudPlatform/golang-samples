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

func TestFilterOnStringValueCustomAttribute(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	tc := testutil.SystemTest(t)
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := filterOnStringValueCustomAttribute(buf, tc.ProjectID); err != nil {
			r.Errorf("filterOnStringValueCustomAttribute: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("filterOnStringValueCustomAttribute got %q, want to contain %q", got, want)
		}
	})
}

func TestFilterOnLongValueCustomAttribute(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	tc := testutil.SystemTest(t)
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := filterOnLongValueCustomAttribute(buf, tc.ProjectID); err != nil {
			r.Errorf("filterOnLongValueCustomAttribute: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("filterOnLongValueCustomAttribute got %q, want to contain %q", got, want)
		}
	})
}

func TestFilterOnMultiCustomAttributes(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	tc := testutil.SystemTest(t)
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := filterOnMultiCustomAttributes(buf, tc.ProjectID); err != nil {
			r.Errorf("filterOnMultiCustomAttributes: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("filterOnMultiCustomAttributes got %q, want to contain %q", got, want)
		}
	})
}
