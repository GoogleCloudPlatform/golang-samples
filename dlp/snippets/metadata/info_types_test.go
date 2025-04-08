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
package metadata

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInfoTypes(t *testing.T) {
	testutil.SystemTest(t)
	tests := []struct {
		language string
		filter   string
		want     string
	}{
		{
			want: "TIME",
		},
		{
			language: "en-US",
			want:     "TIME",
		},
		{
			language: "es",
			want:     "DATE",
		},
		{
			filter: "supported_by=INSPECT",
			want:   "GENDER",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.language, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			err := infoTypes(buf, test.language, test.filter)
			if err != nil {
				t.Errorf("infoTypes(%s, %s) = error %q, want substring %q", test.language, test.filter, err, test.want)
			}
			if got := buf.String(); !strings.Contains(got, test.want) {
				t.Errorf("infoTypes(%s, %s) = %s, want substring %q", test.language, test.filter, got, test.want)
			}
		})
	}
}
