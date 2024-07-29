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

package multimodalaudio

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_summarizeAudio(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)
	location := "us-central1"
	modelName := "gemini-1.5-flash-001"

	err := summarizeAudio(buf, tc.ProjectID, location, modelName)
	if err != nil {
		t.Errorf("Test_summarizeAudio: %v", err.Error())
	}
}

func Test_transcribeAudio(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)
	location := "us-central1"
	modelName := "gemini-1.5-flash-001"

	err := transcribeAudio(buf, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("Test_transcribeAudio: %v", err.Error())
	}

	transcript := buf.String()
	transcriptLowercase := strings.ToLower(transcript)
	// We expect these words pronounced in the podcast to be correctly recognized
	// and transcripted
	for _, word := range []string{
		"pixel",
		"feature",
	} {
		if !strings.Contains(transcriptLowercase, word) {
			t.Errorf("expected the word %q in the transcript of %s", word, "gs://cloud-samples-data/generative-ai/audio/pixel.mp3")
		}
	}
}
