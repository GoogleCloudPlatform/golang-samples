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

package functioncalling

import (
	"io"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_functionCallsChat(t *testing.T) {
	tc := testutil.SystemTest(t)

	w := io.Discard
	location := "us-central1"
	modelName := "gemini-1.0-pro-002"

	err := functionCallsChat(w, tc.ProjectID, location, modelName)
	if err != nil {
		t.Errorf("Test_functionCallsChat: %v", err.Error())
	}
}
