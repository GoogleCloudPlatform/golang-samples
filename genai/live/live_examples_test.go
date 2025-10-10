// Copyright 2025 Google LLC
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

package live

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestLiveGeneration(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)
	t.Run("generate Content in live ground googsearch", func(t *testing.T) {
		buf.Reset()
		err := generateGroundSearchWithTxt(buf)
		if err != nil {
			t.Fatalf("generateGroundSearchWithTxt failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("live Function Call With Text in live", func(t *testing.T) {
		buf.Reset()
		err := generateLiveFuncCallWithTxt(buf)
		if err != nil {
			t.Fatalf("generateLiveFuncCallWithTxt failed: %v", err)
		}

		output := buf.String()
		fmt.Printf("output::%+v", output)
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate structured output with txt", func(t *testing.T) {
		buf.Reset()
		if err := generateStructuredOutputWithTxt(buf); err != nil {
			t.Fatalf("generateStructuredOutputWithTxt failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

}
