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
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDeidentifyDataReplaceWithDictionary(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := deidentifyDataReplaceWithDictionary(&buf, tc.ProjectID, "My name is Alicia Abernathy, and my email address is aabernathy@example.com."); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	want1 := "output: My name is Alicia Abernathy, and my email address is izumi@example.com."
	want2 := "output: My name is Alicia Abernathy, and my email address is alex@example.com."
	if !strings.Contains(got, want1) && !strings.Contains(got, want2) {
		t.Errorf("deidentifyDataReplaceWithDictionary got %q, output does not contains value from dictionary", got)
	}
}
