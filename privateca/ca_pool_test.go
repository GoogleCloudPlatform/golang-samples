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
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestCaPoolTests(t *testing.T) {
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	projectId := tc.ProjectID
	location := "us-central1"
	caPoolId := fmt.Sprintf("test-ca-pool-%v-%v", time.Now().Format("2006-01-02"), r.Int())

	buf := &bytes.Buffer{}

	if err := createCaPool(buf, projectId, location, caPoolId); err != nil {
		t.Errorf("createCaPool got err: %v", err)
	}

	expectedResult := "CA Pool created"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createCaPool got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := deleteCaPool(buf, projectId, location, caPoolId); err != nil {
		t.Errorf("deleteCaPool got err: %v", err)
	}

	expectedResult = "CA Pool deleted"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("deleteCaPool got %q, want %q", got, expectedResult)
	}
}
