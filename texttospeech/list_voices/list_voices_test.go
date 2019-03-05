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
