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

package controlledgeneration

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_controlledGenerationResponseMimeType(t *testing.T) {
	tc := testutil.SystemTest(t)
	w := new(bytes.Buffer)

	location := "us-central1"
	modelName := "gemini-1.5-pro-001"

	err := controlledGenerationResponseMimeType(w, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("controlledGenerationResponseMimeType: %v", err.Error())
	}

	// We explicitly requested a response in JSON, so we're expecting
	// a valid JSON array.
	var array []any
	err = json.Unmarshal(w.Bytes(), &array)
	if err != nil {
		t.Errorf(`could not unmarshal response:
%s

into slice because: %v`, w.Bytes(), err)
	}
}
