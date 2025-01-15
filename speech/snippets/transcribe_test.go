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
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var bucket = "cloud-samples-tests"
var recognitionAudioFile = "../resources/commercial_mono.wav"

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestTranscribeModelSelectionGcs(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := transcribe_model_selection_gcs(&buf); err != nil {
		t.Fatalf("error in transcribe model selection gcs %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Transcript:") {
		t.Errorf("transcribe_model_selection_gcs got %q, expected %q", got, "Transcript:")
	}
}

func TestTranscribeDiarizationBeta(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := transcribe_diarization(&buf); err != nil {
		t.Fatalf("error in transcribe diarization %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Speaker") {
		t.Errorf("transcribe_diarization got %q, expected %q", got, "Speaker")
	}
}

func TestTranscribeDiarizationGcsBeta(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := transcribe_diarization_gcs_beta(&buf); err != nil {
		t.Fatalf("error in transcribe diarization gcs %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Speaker") {
		t.Errorf("transcribe_diarization_gcs_beta got %q, expected %q", got, "Speaker")
	}
}
