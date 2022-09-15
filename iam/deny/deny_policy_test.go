// Copyright 2021 Google LLC
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

package snippets

import (
	"bytes"
	"fmt"
	// "github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"strings"
	"testing"
)

func TestDenyPolicySnippets(t *testing.T) {
	// tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	projectID := "sitalaksmi-deny"
	policyID := "test-p-123"

	if err := createDenyPolicy(buf, projectID, policyID); err != nil {
		t.Fatalf("createDenyPolicy: %v", err)
	}

	expectedResult := fmt.Sprintf("Policy %s created", policyID)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createDenyPolicy got %q, want %q", got, expectedResult)
	}

	buf.Reset()
}
