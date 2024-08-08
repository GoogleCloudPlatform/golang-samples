// Copyright 2020 Google LLC
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
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestContextClasses(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	gcsURI := "gs://cloud-samples-data/speech/commercial_mono.wav"
	err := contextClasses(&buf)
	if err != nil {
		t.Fatalf("contextClasses got err: %v", err)
	}

	want := "Alternative"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("contextClasses(%q): got %q, want %q", gcsURI, got, want)
	}
}
