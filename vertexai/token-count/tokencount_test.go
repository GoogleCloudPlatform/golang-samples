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

func Test_countTokens(t *testing.T) {
	tc := testutil.SystemTest(t)

	prompt := "why is the sky blue?"
	location := "us-central1"
	modelName := "gemini-1.5-flash-001"

	var buf bytes.Buffer
	err := countTokens(&buf, prompt, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("Test_countTokens: %v", err.Error())
	}

	answer := buf.String()
	s := strings.TrimPrefix(answer, "Number of tokens for the prompt: ")
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("Test_countTokens: %v", err.Error())
	}

	// "why is the sky blue?" is expected to account for (more or less) 5 tokens
	// Extremely low or high values would not be correct
	if n <= 1 {
		t.Errorf("Expected more than 1 token, got %d", n)
	}
	if n >= 20 {
		t.Errorf("Expected less than 20 tokens, got %d", n)
	}
}
