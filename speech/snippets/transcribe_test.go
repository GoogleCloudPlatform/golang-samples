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

var bucket_v2 = "cloud-samples-tests"
var gcsVideoPath_v2 = "gs://" + bucket + "/speech/Google_Gnome.wav"
var recognitionAudioFile_v2 = "../resources/commercial_mono.wav"
var gcsDiarizationAudioPath_v2 = "gs://" + bucket + "/speech/commercial_mono.wav"

func TestMain_v2(t *testing.T) {
	exitCode := t.Run()
	os.Exit(exitCode)
}

func TestTranscribeModelSelectionGcs_v2(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := transcribe_model_selection_gcs(&buf, gcsVideoPath, "video"); err != nil {
		t.Fatalf("error in transcribe model selection gcs %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Transcript:") {
		t.Errorf("transcribe_model_selection_gcs got %q, expected %q", got, "Transcript:")
	}
}

func TestTranscribeDiarizationBeta_v2(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := transcribe_diarization(&buf, recognitionAudioFile); err != nil {
		t.Fatalf("error in transcribe diarization %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Speaker") {
		t.Errorf("transcribe_diarization got %q, expected %q", got, "Speaker")
	}
}

func TestTranscribeDiarizationGcsBeta_v2(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := transcribe_diarization_gcs_beta(&buf, gcsDiarizationAudioPath); err != nil {
		t.Fatalf("error in transcribe diarization gcs %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Speaker") {
		t.Errorf("transcribe_diarization_gcs_beta got %q, expected %q", got, "Speaker")
	}
}
