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

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	locationID  = "us"
	processorID = "cf8369d73f92e861"
	filePath    = "testdata/invoice.pdf"
	mimeType    = "application/pdf"
)

func TestProcessDocument(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)

	if err := processDocument(buf, tc.ProjectID, locationID, processorID, filePath, mimeType); err != nil {
		t.Errorf("Quickstart: %v", err)
	}

	got := buf.String()
	if want := "Document Text:"; !strings.Contains(got, want) {
		t.Errorf("processDocument got %q, want %q", got, want)
	}
	if want := "Invoice"; !strings.Contains(got, want) {
		t.Errorf("processDocument got %q, want %q", got, want)
	}
}
