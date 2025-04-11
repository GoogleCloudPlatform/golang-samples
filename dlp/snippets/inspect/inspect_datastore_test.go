// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inspect

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestInspectDatastore(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	writeTestDatastoreFiles(t, tc.ProjectID)
	tests := []struct {
		kind string
		want string
	}{
		{
			kind: "SSNTask",
			want: "Created job",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.kind, func(t *testing.T) {
			t.Parallel()
			testutil.Retry(t, 5, 15*time.Second, func(r *testutil.R) {
				u := uuid.New().String()[:8]
				buf := new(bytes.Buffer)
				if err := inspectDatastore(buf, tc.ProjectID, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, topicName+u, subscriptionName+u, tc.ProjectID, "", test.kind); err != nil {
					r.Errorf("inspectDatastore(%s) got err: %v", test.kind, err)
					return
				}
				if got := buf.String(); !strings.Contains(got, test.want) {
					r.Errorf("inspectDatastore(%s) = %q, want %q substring", test.kind, got, test.want)
				}
			})
		})
	}
}
