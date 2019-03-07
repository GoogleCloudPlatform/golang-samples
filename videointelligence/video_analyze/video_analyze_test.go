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

package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const catVideo = "gs://demomaker/cat.mp4"
const googleworkVideo = "gs://python-docs-samples-tests/video/googlework_short.mp4"

func TestAnalyze(t *testing.T) {
	testutil.SystemTest(t)

	tests := []struct {
		name        string
		gcs         func(io.Writer, string) error
		path        string
		wantContain string
	}{
		{"ShotChange", shotChangeURI, catVideo, "Shot"},
		{"Labels", labelURI, catVideo, "cat"},
		{"Explicit", explicitContentURI, catVideo, "VERY_UNLIKELY"},
		{"SpeechTranscription", speechTranscriptionURI, googleworkVideo, "cultural"},
	}

	for _, tt := range tests {
		if tt.gcs == nil {
			continue
		}

		var buf bytes.Buffer
		err := tt.gcs(&buf, tt.path)
		if err != nil {
			t.Fatalf("GCS %s(%q): got %v, want nil err", tt.name, tt.path, err)
		}
		if got := buf.String(); !strings.Contains(got, tt.wantContain) {
			t.Errorf("GCS %s(%q): got %q, want to contain %q", tt.name, tt.path, got, tt.wantContain)
		}
	}
}
