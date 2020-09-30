// Copyright 2019 Google LLC
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

package v3

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDetectLanguage(t *testing.T) {
	tc := testutil.SystemTest(t)

	text := "Hæ sæta"

	// Detect language.
	var buf bytes.Buffer
	if err := detectLanguage(&buf, tc.ProjectID, text); err != nil {
		t.Fatalf("detectLanguage: %v", err)
	}
	if got, want := buf.String(), "is"; !strings.Contains(got, want) {
		t.Errorf("detectLanguage got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
}
