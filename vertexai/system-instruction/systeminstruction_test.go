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

package systeminstruction

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_systemInstruction(t *testing.T) {
	tc := testutil.SystemTest(t)

	instruction := `
			You are a helpful language translator.
			Your mission is to translate text in English to French.`
	prompt := `
			User input: I like bagels.
    		Answer:`
	location := "us-central1"
	modelName := "gemini-1.0-pro"

	var buf bytes.Buffer
	err := systemInstruction(&buf, instruction, prompt, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("Test_systemInstruction: %v", err.Error())
	}

	answer := buf.String()
	answerLowercase := strings.ToLower(answer)

	// We expect an answer that looks like "J'aime les bagels"
	// The answer being written in French proves that the system instruction was acknowledged
	// by the model
	expected := "J'aime les bagels"
	if !strings.Contains(answerLowercase, "aime") &&
		!strings.Contains(answerLowercase, "adore") &&
		!strings.Contains(answerLowercase, "apprécie") {
		t.Errorf("expected answer like %q, got %q", expected, answer)
	}
	if !strings.Contains(answerLowercase, "bagels") &&
		!strings.Contains(answerLowercase, "baguels") &&
		!strings.Contains(answerLowercase, "beguels") {
		t.Errorf("expected the word %q in answer, got %q", "bagels", answer)
	}
}
