// Copyright 2024 Google LLC
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
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var recognitionAudioFileRawLINEAR16 = "../testdata/audio.raw"

func TestTranscribeStreamingV2SpecificDecoding(t *testing.T) {
	testutil.SystemTest(t)
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	var buf bytes.Buffer
	if err := transcribeStreamingSpecificDecodingV2(&buf, projectID, recognitionAudioFileRawLINEAR16); err != nil {
		t.Fatalf("error in transcribe rawfile %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Brooklyn Bridge") {
		t.Errorf("transcribe_streaming_v2_explicit_decoding got %q, expected %q", got, "Speaker")
	}
}
