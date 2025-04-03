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
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_controlledGenerationResponseSchema6(t *testing.T) {
	tc := testutil.SystemTest(t)
	w := new(bytes.Buffer)

	location := "us-central1"
	modelName := "gemini-1.5-pro-001"

	err := controlledGenerationResponseSchema6(w, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("controlledGenerationResponseSchema6: %v", err.Error())
	}

	// We explicitly requested a response in JSON with a specific schema, so we're
	// expecting to properly decode the output as a slice.
	type Item struct {
		Object string `json:"object"`
	}
	var items [][]Item
	err = json.Unmarshal(w.Bytes(), &items)
	if err != nil {
		t.Errorf(`could not unmarshal response:
%s

into [][]Item because: %v`, w.Bytes(), err)
	}
	if len(items) == 0 {
		t.Errorf("no items returned")
	}
}
