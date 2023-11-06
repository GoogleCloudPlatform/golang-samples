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

func TestDeIdentifyDeterministic(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "Jack's phone number is 5555551212"
	infoTypeNames := []string{"PHONE_NUMBER"}
	keyRingName, err := createKeyRing(t, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	keyFileName, cryptoKeyName, keyVersion, err := createKey(t, tc.ProjectID, keyRingName)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyKey(t, tc.ProjectID, keyVersion)

	surrogateInfoType := "PHONE_TOKEN"
	want := "output : Jack's phone number is PHONE_TOKEN(36):"

	var buf bytes.Buffer

	if err := deIdentifyDeterministicEncryption(&buf, tc.ProjectID, input, infoTypeNames, keyFileName, cryptoKeyName, surrogateInfoType); err != nil {
		t.Errorf("deIdentifyDeterministicEncryption(%q) = error '%q', want %q", err, input, want)
	}

	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("deIdentifyDeterministicEncryption(%q) = %q, want %q", input, got, want)
	}

}
