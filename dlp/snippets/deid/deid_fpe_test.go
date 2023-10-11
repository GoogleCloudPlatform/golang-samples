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

package deid

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDeidTextDataWithFPE(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	input := "My SSN is 123456789"
	infoTypeNames := []string{"US_SOCIAL_SECURITY_NUMBER"}
	surrogateInfoType := "AGE"

	keyRingName, err := createKeyRing(t, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	keyFileName, cryptoKeyName, keyVersion, err := createKey(t, tc.ProjectID, keyRingName)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyKey(t, tc.ProjectID, keyVersion)

	if err := deidentifyFPE(&buf, tc.ProjectID, input, infoTypeNames, keyFileName, cryptoKeyName, surrogateInfoType); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "My SSN is AGE(9): "; strings.Contains(got, want) {
		t.Errorf("deidentifyFPE got %q, want %q", got, want)
	}
}
