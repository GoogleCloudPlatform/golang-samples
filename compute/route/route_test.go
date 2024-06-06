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
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateRoute(t *testing.T) {
	buf := &bytes.Buffer{}
	tc := testutil.SystemTest(t)
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	routeName := fmt.Sprintf("testname-%v", r.Int())

	if err := createRoute(buf, tc.ProjectID, routeName); err != nil {
		t.Errorf("createRoute got error: %v", err)
	}

	if err := listRoutes(buf, tc.ProjectID); err != nil {
		t.Errorf("listRoutes got error: %v", err)
	}

	expectedResult := fmt.Sprintf("- %s", routeName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("listInstances got %q, want %q", got, expectedResult)
	}

	if err := deleteRoute(buf, tc.ProjectID, routeName); err != nil {
		t.Errorf("deleteRoute got error: %v", err)
	}

	expectedResult2 := "Route deleted"
	if got := buf.String(); !strings.Contains(got, expectedResult2) {
		t.Errorf("deleteRoute got %q, want %q", got, expectedResult)
	}
	buf.Reset()
}
