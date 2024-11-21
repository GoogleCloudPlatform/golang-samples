// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package embeddings

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestEmbeddings(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	location := "us-central1"

	t.Run("generate embeddings with lower dimension", func(t *testing.T) {
		buf.Reset()
		err := generateWithLowerDimension(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateWithLowerDimension failed: %v", err)
		}

		expectedOutputs := []string{
			"Text embedding (length=128): [",
			"Image embedding (length=128): [",
		}
		actOutput := buf.String()
		for _, exp := range expectedOutputs {
			if !strings.Contains(actOutput, exp) {
				t.Errorf("expected output to contain text %q, got: %q", exp, actOutput)
			}
		}
	})

	t.Run("generate embeddings for image and text", func(t *testing.T) {
		buf.Reset()
		err := generateForTextAndImage(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateForImageAndText failed: %v", err)
		}

		expectedOutputs := []string{
			"Text embedding (length=1408): [",
			"Image embedding (length=1408): [",
		}
		actOutput := buf.String()
		for _, exp := range expectedOutputs {
			if !strings.Contains(actOutput, exp) {
				t.Errorf("expected output to contain text %q, got: %q", exp, actOutput)
			}
		}
	})

	t.Run("generate embedding for image", func(t *testing.T) {
		buf.Reset()
		err := generateForImage(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateForImage failed: %v", err)
		}

		expOutput := "Image embedding (length=1408): ["
		actOutput := buf.String()
		if !strings.Contains(actOutput, expOutput) {
			t.Errorf("expected output to contain text %q, got: %q", expOutput, actOutput)
		}
	})

	t.Run("generate embedding for video", func(t *testing.T) {
		buf.Reset()
		err := generateForVideo(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateForVideo failed: %v", err)
		}

		expOutput := "Video embedding (seconds: 1-5; length=1408): ["
		actOutput := buf.String()
		if !strings.Contains(actOutput, expOutput) {
			t.Errorf("expected output to contain text %q, got: %q", expOutput, actOutput)
		}
	})

	t.Run("generate embeddings for image text and video", func(t *testing.T) {
		buf.Reset()
		err := generateForImageTextAndVideo(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateForImageTextAndVideo failed: %v", err)
		}

		expectedOutputs := []string{
			"Text embedding (length=1408): [",
			"Image embedding (length=1408): [",
			"Video embedding (length=1408): [",
		}
		actOutput := buf.String()
		for _, exp := range expectedOutputs {
			if !strings.Contains(actOutput, exp) {
				t.Errorf("expected output to contain text %q, got: %q", exp, actOutput)
			}
		}
	})
}
