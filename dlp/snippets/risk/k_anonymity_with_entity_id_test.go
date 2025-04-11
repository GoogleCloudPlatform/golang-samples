// Copyright 2023 Google LLC
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

package risk

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCalculateKAnonymityWithEntityId(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)
	err := calculateKAnonymityWithEntityId(buf, tc.ProjectID, dataSetID, tableID, "title", "contributor_ip")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Created job"; !strings.Contains(got, want) {
		t.Errorf("CalculateKAnonymityWithEntityId got %s, want substring %q", got, want)
	}
	if got, want := buf.String(), "Quasi-ID values"; !strings.Contains(got, want) {
		t.Errorf("CalculateKAnonymityWithEntityId got %s, want substring %q", got, want)
	}
	if got, want := buf.String(), "Job status: DONE"; !strings.Contains(got, want) {
		t.Errorf("CalculateKAnonymityWithEntityId got %s, want substring %q", got, want)
	}
	if got, want := buf.String(), "Bucket size range:"; !strings.Contains(got, want) {
		t.Errorf("CalculateKAnonymityWithEntityId got %s, want substring %q", got, want)
	}
	if got, want := buf.String(), "Class size: 1"; !strings.Contains(got, want) {
		t.Errorf("CalculateKAnonymityWithEntityId got %s, want substring %q", got, want)
	}
	if got, want := buf.String(), "Job name:"; !strings.Contains(got, want) {
		t.Errorf("CalculateKAnonymityWithEntityId got %s, want substring %q", got, want)
	}

}
