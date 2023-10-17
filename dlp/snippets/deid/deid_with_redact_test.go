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
package deid

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDeIdentifyWithRedact(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "My name is Alicia Abernathy, and my email address is aabernathy@example.com."
	infoTypeNames := []string{"EMAIL_ADDRESS"}
	want := "output: My name is Alicia Abernathy, and my email address is ."

	var buf bytes.Buffer

	if err := deidentifyWithRedact(&buf, tc.ProjectID, input, infoTypeNames); err != nil {
		t.Errorf("deidentifyWithRedact(%q) = error '%q', want %q", err, input, want)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyWithRedact(%q) = %q, want %q", got, input, want)
	}
}
