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
package redact

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRedactImageFileAllText(t *testing.T) {
	tc := testutil.SystemTest(t)
	inputPath := "testdata/image.jpg"
	outputPath := "testdata/test-output-sensitive-data-image-redacted.jpeg"

	var buf bytes.Buffer
	if err := redactImageFileAllText(&buf, tc.ProjectID, inputPath, outputPath); err != nil {
		t.Fatal(err)
	}

	hash1, err := calculateImageHash(inputPath)
	if err != nil {
		t.Errorf("redactImageFileAllText: Error calculating hash for image 1: %q", err)
	}

	if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
		t.Error("redactImageFileAllText: the output file is not generated")
	} else {
		hash2, err := calculateImageHash(outputPath)
		if err != nil {
			t.Errorf("redactImageFileAllText: Error calculating hash for image 2: %q", err)
		}

		if hash1 == hash2 {
			t.Error("redactImageFileAllText: image is not redacted.")
		}
	}

	got := buf.String()
	if want := "Wrote output to"; !strings.Contains(got, want) {
		t.Errorf("redactImageFileAllText got %q, want %q", got, want)
	}

}
