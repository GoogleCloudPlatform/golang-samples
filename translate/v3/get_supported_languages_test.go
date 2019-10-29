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

func TestGetSupportedLanguages(t *testing.T) {
	tc := testutil.SystemTest(t)

	// Get supported languages.
	var buf bytes.Buffer
	if err := getSupportedLanguages(&buf, tc.ProjectID); err != nil {
		t.Fatalf("getSupportedLanguages: %v", err)
	}
	if got, want := buf.String(), "zh-CN"; !strings.Contains(got, want) {
		t.Errorf("getSupportedLanguages got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
}
