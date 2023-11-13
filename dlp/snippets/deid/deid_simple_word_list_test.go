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

func TestDeIdentifyWithWordList(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	input := "Patient was seen in RM-YELLOW then transferred to rm green."
	infoType := "CUSTOM_ROOM_ID"
	wordList := []string{"RM-GREEN", "RM-YELLOW", "RM-ORANGE"}
	want := "output : Patient was seen in [CUSTOM_ROOM_ID] then transferred to [CUSTOM_ROOM_ID]."

	if err := deidentifyWithWordList(&buf, tc.ProjectID, input, infoType, wordList); err != nil {
		t.Errorf("deidentifyWithWordList(%q) = error '%q', want %q", input, err, want)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyWithWordList(%q) = %q, want %q", input, got, want)
	}
}
