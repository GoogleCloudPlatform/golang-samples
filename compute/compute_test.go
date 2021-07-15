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
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestComputeSnippets(t *testing.T) {
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	buf := &bytes.Buffer{}

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := listInstances(buf, tc.ProjectID, zone); err != nil {
			r.Errorf("listInstances got err: %v", err)
		}

		expectedResult := "Instances found in zone"
		if got := buf.String(); !strings.Contains(got, "Instances found in zone") {
			r.Errorf("listInstances got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, expectedResult)
		}
	})

}
