// Copyright 2024 Google LLC
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
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateRoute(t *testing.T) {
	buf := &bytes.Buffer{}
	tc := testutil.SystemTest(t)
	name := "testname"

	if err := createRoute(buf, tc.ProjectID, name); err != nil {
		t.Fatalf("createRoute got error: %v", err)
	}

	buf.Reset()

	if err := listRoutes(buf, tc.ProjectID); err != nil {
		t.Fatalf("listRoutes got error: %v", err)
	}

	expectedResult := fmt.Sprintf("- %s", "testname")
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("listInstances got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := deleteRoute(buf, tc.ProjectID, name); err != nil {
		t.Fatalf("deleteRoute got error: %v", err)
	}

	expectedResult = "Route deleted"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("deleteRoute got %q, want %q", got, expectedResult)
	}

	buf.Reset()
}
