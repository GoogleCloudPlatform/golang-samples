// Copyright 2023 Google LLC
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
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestPolicy(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	policy, err := getPolicy(buf, tc.ProjectID)
	if err != nil {
		t.Fatalf("getPolicy(%q) failed: %v", tc.ProjectID, err)
	}
	want := "bindings"
	if !strings.Contains(buf.String(), want) {
		t.Errorf("getPolicy got %q, want to contain %q", buf.String(), want)
	}

	buf.Reset()
	err = setPolicy(buf, tc.ProjectID, policy)
	if err != nil {
		t.Fatalf("setPolicy(%q, %v) failed: %v", tc.ProjectID, policy, err)
	}
	if !strings.Contains(buf.String(), want) {
		t.Errorf("setPolicy got %q, want to contain %q", buf.String(), want)
	}
}
