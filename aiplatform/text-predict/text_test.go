// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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

func TestTextPredict(t *testing.T) {
	tc := testutil.SystemTest(t)

	prompt := "hello, say something nice."
	projectID := tc.ProjectID
	location := "us-central1"
	publisher := "google"
	model := "text-bison"

	parameters := map[string]interface{}{
		"temperature":     0.8,
		"maxOutputTokens": 256,
		"topP":            0.4,
		"topK":            40,
	}

	var buf bytes.Buffer
	if err := textPredict(&buf, prompt, projectID, location, publisher, model, parameters); err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, "text-prediction response") {
		t.Error("generated text content not found in response")
	}
}
