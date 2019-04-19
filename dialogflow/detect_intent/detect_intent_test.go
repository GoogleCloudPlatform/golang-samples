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
	"fmt"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDetectIntentText(t *testing.T) {
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID

	sessionID := fmt.Sprintf("golang-samples-test-session-%v", time.Now())

	text := "I'd like to book a room"

	languageCode := "en-US"

	_, err := DetectIntentText(projectID, sessionID, text, languageCode)

	if err != nil {
		t.Error(err)
	}
}

func TestDetectIntentAudio(t *testing.T) {
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID

	sessionID := fmt.Sprintf("golang-samples-test-session-%v", time.Now())

	audioFile := "../resources/book_a_room.wav"

	languageCode := "en-US"

	_, err := DetectIntentAudio(projectID, sessionID, audioFile, languageCode)

	if err != nil {
		t.Error(err)
	}
}

func TestDetectIntentAudioWithNonexistentFile(t *testing.T) {
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID

	sessionID := fmt.Sprintf("golang-samples-test-session-%v", time.Now())

	audioFile := "./this-file-should-not-exist.wav"

	languageCode := "en-US"

	_, err := DetectIntentAudio(projectID, sessionID, audioFile, languageCode)

	if err == nil {
		t.Error("Expected due to non-existent file")
	}
}

func TestDetectIntentStream(t *testing.T) {
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID

	sessionID := fmt.Sprintf("golang-samples-test-session-%v", time.Now())

	audioFile := "../resources/book_a_room.wav"

	languageCode := "en-US"

	_, err := DetectIntentAudio(projectID, sessionID, audioFile, languageCode)

	if err != nil {
		t.Error(err)
	}

}

func TestDetectIntentStreamWithNonexistentFile(t *testing.T) {
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID

	sessionID := fmt.Sprintf("golang-samples-test-session-%v", time.Now())

	audioFile := "./this-file-should-not-exist.wav"

	languageCode := "en-US"

	_, err := DetectIntentStream(projectID, sessionID, audioFile, languageCode)

	if err == nil {
		t.Error("Expected due to non-existent file")
	}
}
