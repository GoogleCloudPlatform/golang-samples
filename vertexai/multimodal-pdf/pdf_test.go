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

package multimodalpdf

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_generateContentFromPDF(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)
	prompt := pdfPrompt{
		pdfPath: "gs://cloud-samples-data/generative-ai/pdf/2403.05530.pdf",
		question: `
			You are a very professional document summarization specialist.
    		Please summarize the given document.
		`,
	}
	location := "us-central1"
	modelName := "gemini-1.5-pro-preview-0409"

	err := generateContentFromPDF(buf, prompt, tc.ProjectID, location, modelName)
	if err != nil {
		t.Errorf("Test_generateContentFromPDF: %v", err.Error())
	}

	generatedSummary := buf.String()
	generatedSummaryLowercase := strings.ToLower(generatedSummary)
	// We expect these important topics in the video to be correctly covered
	// in the generated summary
	for _, word := range []string{
		"gemini",
		"tokens",
	} {
		if !strings.Contains(generatedSummaryLowercase, word) {
			t.Errorf("expected the word %q in the description of %s", word, prompt.pdfPath)
		}
	}
}
