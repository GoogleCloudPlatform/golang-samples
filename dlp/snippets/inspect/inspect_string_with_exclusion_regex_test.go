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

func TestInspectStringWithExclusionRegex(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	if err := inspectStringWithExclusionRegex(&buf, tc.ProjectID, "Some email addresses: gary@example.com, bob@example.org", ".+@example.com"); err != nil {
		t.Errorf("inspectStringWithExclusionRegex: %v", err)
	}

	got := buf.String()

	if want := "Quote: bob@example.org"; !strings.Contains(got, want) {
		t.Errorf("inspectStringWithExclusionRegex got %q, want %q", got, want)
	}
	if want := "Quote: gary@example.com"; strings.Contains(got, want) {
		t.Errorf("inspectStringWithExclusionRegex got %q, want %q", got, want)
	}
}
