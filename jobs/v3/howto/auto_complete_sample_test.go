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

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestJobTitleAutoComplete(t *testing.T) {
	t.Skip("Flaky. https://github.com/GoogleCloudPlatform/golang-samples/issues/1061.")

	tc := testutil.SystemTest(t)

	buf := &bytes.Buffer{}
	if _, err := jobTitleAutoComplete(buf, tc.ProjectID, "", "SWE"); err != nil {
		t.Fatalf("jobTitleAutoComplete: %v", err)
	}
	want := "Auto complete results"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("jobTitleAutoComplete got %q, want %q", got, want)
	}

	buf.Reset()
	if _, err := defaultAutoComplete(buf, tc.ProjectID, "", "SWE"); err != nil {
		t.Fatalf("defaultAutoComplete: %v", err)
	}
	want = "Auto complete results"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("defaultAutoComplete got %q, want %q", got, want)
	}
}

func TestDefaultAutoComplete(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := &bytes.Buffer{}
	if _, err := defaultAutoComplete(buf, tc.ProjectID, "", "SWE"); err != nil {
		t.Fatalf("defaultAutoComplete: %v", err)
	}
	want := "Auto complete results"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("defaultAutoComplete got %q, want %q", got, want)
	}
}
