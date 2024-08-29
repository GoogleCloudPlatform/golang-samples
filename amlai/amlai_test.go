// Copyright 2024 Google LLC
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

package amlai

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// To run the tests, do the following:
// Export the following env vars:
// *   GOOGLE_APPLICATION_CREDENTIALS
// *   GOLANG_SAMPLES_PROJECT_ID
// Enable the following API on the test project:
// *   AML AI API

// TestLocations tests the GET operation on locations.
func TestLocations(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	location := "us-central1"

	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := listLocations(buf, tc.ProjectID); err != nil {
			r.Errorf("listLocations got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, location) {
			r.Errorf("listLocations got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, location)
		}
	})
	buf.Reset()
}
