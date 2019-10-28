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

package v3

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestTranslateTextWithGlossaryAndModel(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	var got string

	location := "us-central1"
	sourceLang := "en"
	targetLang := "ja"
	text := "That' il do it. deception"
	glossaryID := fmt.Sprintf("translate_text_with_glossary_and_model-%v", uuid.New().ID())
	glossaryInputURI := "gs://cloud-samples-data/translation/glossary_ja.csv"
	modelID := "TRL3128559826197068699"

	// Create a glossary.
	if err := createGlossary(&buf, tc.ProjectID, location, glossaryID, glossaryInputURI); err != nil {
		t.Fatalf("createGlossary: %v", err)
	}
	defer deleteGlossary(&buf, tc.ProjectID, location, glossaryID)

	// Translate text.
	if err := translateTextWithGlossaryAndModel(&buf, tc.ProjectID, location, sourceLang, targetLang, text, glossaryID, modelID); err != nil {
		t.Fatalf("translateTextWithGlossaryAndModel: %v", err)
	}

	got = buf.String()

	// Custom model.
	if !strings.Contains(got, "それはそうだ") && !strings.Contains(got, "それじゃあ") {
		t.Fatalf("Got '%s', expected to contain 'それはそうだ' or 'それじゃあ'", got)
	}

	// Glossary.
	if !strings.Contains(got, "欺く") {
		t.Fatalf("Got '%s', expected to contain '欺く'", got)
	}
}
