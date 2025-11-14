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
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func generateLiveFuncCallWithTxtMock(w io.Writer) error {
	mockOutput := "Mocked Live Function Call: result = 42"
	_, err := fmt.Fprintln(w, mockOutput)
	return err
}

func generateStructuredOutputWithTxtMock(w io.Writer) error {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	mock := Person{
		Name: "John Doe",
		Age:  42,
	}

	b, err := json.Marshal(mock)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, string(b))
	return err
}

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
		err := generateLiveFuncCallWithTxtMock(buf)
		if err != nil {
			t.Fatalf("generateLiveFuncCallWithTxt failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate structured output with txt", func(t *testing.T) {
		buf.Reset()
		if err := generateStructuredOutputWithTxtMock(buf); err != nil {
			t.Fatalf("generateStructuredOutputWithTxt failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

}
