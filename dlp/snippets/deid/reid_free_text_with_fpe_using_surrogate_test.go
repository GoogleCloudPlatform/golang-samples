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

func TestReidentifyFreeTextWithFPEUsingSurrogate(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	inputStr := "My phone number is 1234567890"
	infoType := "PHONE_NUMBER"
	surrogateType := "PHONE_TOKEN"
	unWrappedKey, err := getUnwrappedKey(t)
	if err != nil {
		t.Fatal(err)
	}

	if err := deidentifyFreeTextWithFPEUsingSurrogate(&buf, tc.ProjectID, inputStr, infoType, surrogateType, unWrappedKey); err != nil {
		t.Fatal(err)
	}

	inputForReid := buf.String()
	inputForReid = strings.TrimPrefix(inputForReid, "output: ")
	buf.Reset()
	if err := reidentifyFreeTextWithFPEUsingSurrogate(&buf, tc.ProjectID, inputForReid, surrogateType, unWrappedKey); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "output: My phone number is 1234567890"; got != want {
		t.Errorf("reidentifyFreeTextWithFPEUsingSurrogate got %q, want %q", got, want)
	}

}
