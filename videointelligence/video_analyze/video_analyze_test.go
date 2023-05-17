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
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const catVideo = "gs://cloud-samples-data/video/cat.mp4"
const googleworkVideo = "gs://python-docs-samples-tests/video/googlework_short.mp4"
const testFailureFormat = "%s failed: wanted %s, got %s"

func assert(buf bytes.Buffer, want string, t *testing.T) {
	got := buf.String()
	if !strings.Contains(got, want) {
		t.Errorf(testFailureFormat, t.Name(), want, got)
	}
}

func TestAnalyzeShotChange(t *testing.T) {
	t.Skip("see GoogleCloudPlatform/golang-samples#3049")
	testutil.EndToEndTest(t)

	testutil.Retry(t, 10, time.Minute, func(r *testutil.R) {
		want := "Shot"
		var buf bytes.Buffer
		err := shotChangeURI(&buf, catVideo)

		if err != nil {
			r.Errorf("%v", err)
		}
		assert(buf, want, t)
	})
}

func TestAnalyzeLabelURI(t *testing.T) {
	t.Skip("see GoogleCloudPlatform/golang-samples#3049")
	testutil.EndToEndTest(t)

	testutil.Retry(t, 10, time.Minute, func(r *testutil.R) {
		want := "cat"
		var buf bytes.Buffer
		err := labelURI(&buf, catVideo)
		if err != nil {
			r.Errorf("%v", err)
		}
		assert(buf, want, t)
	})
}

func TestAnalyzeExplicitContentURI(t *testing.T) {
	t.Skip("see GoogleCloudPlatform/golang-samples#3049")
	testutil.EndToEndTest(t)

	testutil.Retry(t, 10, time.Minute, func(r *testutil.R) {
		want := "VERY_UNLIKELY"
		var buf bytes.Buffer
		err := explicitContentURI(&buf, catVideo)
		if err != nil {
			r.Errorf("%v", err)
		}
		assert(buf, want, t)
	})
}

func TestAnalyzeSpeechTranscriptionURI(t *testing.T) {
	t.Skip("see GoogleCloudPlatform/golang-samples#3049")
	testutil.EndToEndTest(t)

	testutil.Retry(t, 10, time.Minute, func(r *testutil.R) {
		want := "cultural"
		var buf bytes.Buffer
		err := speechTranscriptionURI(&buf, googleworkVideo)
		if err != nil {
			r.Errorf("%v", err)
		}
		assert(buf, want, t)
	})
}
