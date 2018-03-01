// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
