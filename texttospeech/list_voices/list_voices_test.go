// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestListVoices(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	err := ListVoices(&buf)
	if err != nil {
		t.Error(err)
	}
	got := buf.String()

	if !strings.Contains(got, "en-US") {
		t.Error("'en-US' not found")
	}

	if !strings.Contains(got, "SSML Voice Gender: MALE") {
		t.Error("'SSML Voice Gender: MALE' not found")
	}

	if !strings.Contains(got, "SSML Voice Gender: FEMALE") {
		t.Error("'SSML Voice Gender: FEMALE' not found")
	}
}
