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

func TestInspectWithCustomRegex(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := inspectWithCustomRegex(&buf, tc.ProjectID, "Patients MRN 444-5-22222", "[1-9]{3}-[1-9]{1}-[1-9]{5}", "C_MRN"); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Infotype Name: C_MRN"; !strings.Contains(got, want) {
		t.Errorf("inspectWithCustomRegex got %q, want %q", got, want)
	}
	if want := "Likelihood: POSSIBLE"; !strings.Contains(got, want) {
		t.Errorf("inspectWithCustomRegex got %q, want %q", got, want)
	}
}
