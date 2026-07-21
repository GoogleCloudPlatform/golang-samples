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

package tokencount

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_countTokensMultimodal(t *testing.T) {
	t.Skip("Skipped while waiting to decommission vertexai")
	tc := testutil.SystemTest(t)

	location := "us-central1"
	modelName := "gemini-2.0-flash-001"

	var buf bytes.Buffer
	err := countTokensMultimodal(&buf, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("Test_countTokensMultimodal: %v", err.Error())
	}

	answer := buf.String()

	for _, expected := range []string{
		"Number of tokens for the multimodal video prompt: ",
		"Prompt Token Count:",
		"Candidates Token Count:",
		"Total Token Count:",
	} {
		if !strings.Contains(answer, expected) {
			t.Errorf("Response does not contain %q", expected)
		}
	}

	s := strings.TrimPrefix(answer, "Number of tokens for the multimodal video prompt: ")
	s, _, _ = strings.Cut(s, "\n")
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("Test_countTokensMultimodal: %v", err.Error())
	}

	// The pixel8.mp4 video prompt is expected to account for about 17,000 tokens
	// Extremely low or high values would not be correct
	if n <= 100 {
		t.Errorf("Expected more than 100 tokens, got %d", n)
	}
	if n >= 1_000_000 {
		t.Errorf("Expected less than 1,000,000 tokens, got %d", n)
	}
}
