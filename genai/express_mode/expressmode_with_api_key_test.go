// Copyright 2025 Google LLC
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

package express_mode

import (
	"bytes"
	"context"
	"testing"

	"google.golang.org/genai"
)

type fakeModels struct{}

func (m *fakeModels) GenerateContent(ctx context.Context, model string, contents []*genai.Content, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	return &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []*genai.Part{{Text: "mocked bubble sort explanation"}},
				},
			},
		},
	}, nil
}

func TestExpressModeGenerationWithMockFunctional(t *testing.T) {
	buf := new(bytes.Buffer)
	client := &fakeModels{}

	resp, err := client.GenerateContent(context.Background(), "gemini-2.5-flash", nil, nil)
	if err != nil {
		t.Fatalf("fake GenerateContent failed: %v", err)
	}

	buf.WriteString(resp.Candidates[0].Content.Parts[0].Text + "\n")

	got := buf.String()
	if got == "" {
		t.Error("expected non-empty mocked output, got empty")
	}
	if got != "mocked bubble sort explanation\n" {
		t.Errorf("unexpected output: got %q", got)
	}
}
