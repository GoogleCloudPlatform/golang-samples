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

func TestGetSink(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	sinkName := "_Default"
	if err := getSink(buf, tc.ProjectID, sinkName); err != nil {
		t.Fatalf("getSink(%q, %q) failed: %v", tc.ProjectID, sinkName, err)
	}
	if !strings.Contains(buf.String(), sinkName) {
		t.Errorf("getSink got %q, want to contain %q", buf.String(), sinkName)
	}
}
