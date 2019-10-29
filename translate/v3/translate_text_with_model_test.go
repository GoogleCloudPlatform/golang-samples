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

// TODO: uncomment this test once the AutoML model is set in the testing projects.

// import (
// 	"bytes"
// 	"strings"
// 	"testing"

// 	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
// )

// func TestTranslateTextWithModel(t *testing.T) {
// 	tc := testutil.SystemTest(t)

// 	location := "us-central1"
// 	sourceLang := "en"
// 	targetLang := "ja"
// 	text := "That' il do it."
// 	modelID := "TRL3128559826197068699"

// 	// Translate text.
// 	var buf bytes.Buffer
// 	if err := translateTextWithModel(&buf, tc.ProjectID, location, sourceLang, targetLang, text, modelID); err != nil {
// 		t.Fatalf("translateTextWithModel: %v", err)
// 	}
// 	if got, want1, want2 := buf.String(), "それはそうだ", "それじゃあ"; !strings.Contains(got, want1) && !strings.Contains(got, want2) {
// 		t.Errorf("translateTextWithModel got:\n----\n%s----\nWant to contain:\n----\n%s\n----\nOR\n----\n%s\n----", got, want1, want2)
// 	}
// }
