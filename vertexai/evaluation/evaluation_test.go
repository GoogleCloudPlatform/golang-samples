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

package evaluation

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestEvaluation(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)
	location := "us-central1"

	err := getROUGEScore(buf, tc.ProjectID, location)
	if err != nil {
		t.Fatalf("getRougeScore: %v", err.Error())
	}

	buf.Reset()
	err = evaluateOutput(buf, tc.ProjectID, location)
	if err != nil {
		t.Fatalf("evaluateOutput: %v", err.Error())
	}

	buf.Reset()
	err = pairwiseEvaluation(buf, tc.ProjectID, location)
	if err != nil {
		t.Fatalf("pairwiseEvaluation: %v", err.Error())
	}
}
