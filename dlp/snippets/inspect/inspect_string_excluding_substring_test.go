// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package inspect

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInspectStringCustomExcludingSubstring(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectStringCustomExcludingSubstring(&buf, tc.ProjectID, "Name: Doe, John. Name: Example, Jimmy", "[A-Z][a-z]{1,15}, [A-Z][a-z]{1,15}", []string{"Jimmy"}); err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	if want := "Infotype Name: CUSTOM_NAME_DETECTOR"; !strings.Contains(got, want) {
		t.Errorf("inspectStringCustomExcludingSubstring got %q, want %q", got, want)
	}
	if want := "Quote: Doe, John"; !strings.Contains(got, want) {
		t.Errorf("inspectStringCustomExcludingSubstring got %q, want %q", got, want)
	}
	if want := "Jimmy"; strings.Contains(got, want) {
		t.Errorf("inspectStringCustomExcludingSubstring got %q, want %q", got, want)
	}
}
