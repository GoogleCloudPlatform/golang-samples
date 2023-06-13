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

package redact

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/dlp/apiv2/dlppb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRedactImage(t *testing.T) {
	tc := testutil.SystemTest(t)
	tests := []struct {
		name      string
		inputPath string
		bt        dlppb.ByteContentItem_BytesType
		infoTypes []string
		want      string
	}{
		{
			name:      "image with one type",
			inputPath: "testdata/ok.png",
			bt:        dlppb.ByteContentItem_IMAGE_PNG,
			infoTypes: []string{"US_SOCIAL_SECURITY_NUMBER"},
			want:      "Wrote output to",
		},
		{
			name:      "image with two types",
			inputPath: "testdata/ok.png",
			bt:        dlppb.ByteContentItem_IMAGE_PNG,
			infoTypes: []string{"US_SOCIAL_SECURITY_NUMBER", "DATE"},
			want:      "Wrote output to",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			// TODO: output to a Writer or bytes rather than to a file on disk.
			if err := redactImage(buf, tc.ProjectID, test.infoTypes, test.bt, test.inputPath, "testdata/test_output.png"); err != nil {
				t.Errorf("redactImage: %v", err)
			}
			if got := buf.String(); !strings.Contains(got, test.want) {
				t.Errorf("redactImage(%s) got %q, want substring %q", test.name, got, test.want)
			}
		})
	}
}

func TestRedactImageFileListedInfoTypes(t *testing.T) {
	tc := testutil.SystemTest(t)
	inputPath := "testdata/image.jpg"
	outputPath := "testdata/test-output-image-file-listed-infoTypes-redacted.jpeg"

	var buf bytes.Buffer
	if err := redactImageFileListedInfoTypes(&buf, tc.ProjectID, inputPath, outputPath); err != nil {
		t.Fatal(err)
	}

	hash1, err := calculateImageHash(inputPath)
	if err != nil {
		t.Errorf("redactImageFileListedInfoTypes: Error calculating hash for image 1: %q", err)
	}

	if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
		t.Error("redactImageFileListedInfoTypes: the output file is not generated")
	} else {
		hash2, err := calculateImageHash(outputPath)
		if err != nil {
			t.Errorf("redactImageFileListedInfoTypes: Error calculating hash for image 2: %q", err)
		}

		if hash1 == hash2 {
			t.Error("redactImageFileListedInfoTypes: image is not redacted.")
		}
	}

	got := buf.String()
	if want := "Wrote output to"; !strings.Contains(got, want) {
		t.Errorf("redactImageFileListedInfoTypes got %q, want %q", got, want)
	}

	if want := "ioutil.ReadFile: open testdata/image.jpg: The system cannot find the path specified."; strings.Contains(got, want) {
		t.Errorf("redactImageFileListedInfoTypes got %q, want %q", got, want)
	}

}

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

func TestRedactImageFileColoredInfoTypes(t *testing.T) {
	tc := testutil.SystemTest(t)

	inputPath := "testdata/image.jpg"
	outputPath := "testdata/test-output-image-file-colored-infoTypes-redacted.jpeg"

	var buf bytes.Buffer
	if err := redactImageFileColoredInfoTypes(&buf, tc.ProjectID, inputPath, outputPath); err != nil {
		t.Fatal(err)
	}

	hash1, err := calculateImageHash(inputPath)
	if err != nil {
		t.Errorf("redactImageFileColoredInfoTypes: Error calculating hash for image 1: %q", err)
	}

	if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
		t.Error("redactImageFileColoredInfoTypes: the output file is not generated")
	} else {
		hash2, err := calculateImageHash(outputPath)
		if err != nil {
			t.Errorf("redactImageFileColoredInfoTypes: Error calculating hash for image 2: %q", err)
		}

		if hash1 == hash2 {
			t.Error("redactImageFileColoredInfoTypes: image is not redacted.")
		}
	}

	got := buf.String()
	if want := "Wrote output to"; !strings.Contains(got, want) {
		t.Errorf("redactImageFileColoredInfoTypes got %q, want %q", got, want)
	}

	if want := "ioutil.ReadFile: open testdata/image.jpg: The system cannot find the path specified."; strings.Contains(got, want) {
		t.Errorf("redactImageFileColoredInfoTypes got %q, want %q", got, want)
	}

}

func calculateImageHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashSum := hash.Sum(nil)
	return hex.EncodeToString(hashSum), nil
}
