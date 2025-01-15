// Copyright 2024 Google LLC
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

package multimodalvideoaudio

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_generateMultimodalContent(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)
	location := "us-central1"
	modelName := "gemini-1.5-flash-001"

	err := generateMultimodalContent(buf, tc.ProjectID, location, modelName)
	if err != nil {
		t.Errorf("Test_generateMultimodalContent: %v", err.Error())
	}

	generatedDescription := buf.String()
	generatedDescriptionLowercase := strings.ToLower(generatedDescription)
	// We expect these important topics in the video to be correctly covered
	// in the generated description
	for _, word := range []string{
		"tokyo",
		"pixel",
	} {
		if !strings.Contains(generatedDescriptionLowercase, word) {
			t.Errorf("expected the word %q in the description of the video", word)
		}
	}
}
