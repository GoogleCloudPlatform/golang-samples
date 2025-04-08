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

package deid

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDeidentifyDateShift(t *testing.T) {
	tc := testutil.SystemTest(t)
	tests := []struct {
		input      string
		want       string
		lowerBound int32
		upperBound int32
	}{
		{
			input:      "2016-01-10",
			lowerBound: 1,
			upperBound: 1,
			want:       "2016-01-11",
		},
		{
			input:      "2016-01-10",
			lowerBound: -1,
			upperBound: -1,
			want:       "2016-01-09",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			test := test
			t.Parallel()
			buf := new(bytes.Buffer)
			err := deidentifyDateShift(buf, tc.ProjectID, test.lowerBound, test.upperBound, test.input)
			if err != nil {
				t.Errorf("deidentifyDateShift(%v, %v, %q) = error '%q', want %q", test.lowerBound, test.upperBound, err, test.input, test.want)
			}
			if got := buf.String(); got != test.want {
				t.Errorf("deidentifyDateShift(%v, %v, %q) = %q, want %q", test.lowerBound, test.upperBound, got, test.input, test.want)
			}
		})
	}
}
