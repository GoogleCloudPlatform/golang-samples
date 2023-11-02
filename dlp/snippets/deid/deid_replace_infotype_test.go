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

func TestDeIdentifyWithInfotype(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "My email is test@example.com"
	infoType := []string{"EMAIL_ADDRESS"}
	want := "output : My email is [EMAIL_ADDRESS]"

	var buf bytes.Buffer

	if err := deidentifyWithInfotype(&buf, tc.ProjectID, input, infoType); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyFreeTextWithFPEUsingSurrogate(%q) = %q, want %q", input, got, want)
	}

}
