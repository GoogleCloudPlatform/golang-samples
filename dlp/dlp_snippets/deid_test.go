// Copyright 2019 Google LLC
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

package main

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMask(t *testing.T) {
	testutil.SystemTest(t)
	tests := []struct {
		input            string
		maskingCharacter string
		numberToMask     int32
		want             string
	}{
		{
			input:            "My SSN is 111222333",
			maskingCharacter: "+",
			want:             "My SSN is +++++++++",
		},
		{
			input: "My SSN is 111222333",
			want:  "My SSN is *********",
		},
		{
			input:            "My SSN is 111222333",
			maskingCharacter: "+",
			numberToMask:     6,
			want:             "My SSN is ++++++333",
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			mask(buf, client, projectID, test.input, []string{"US_SOCIAL_SECURITY_NUMBER"}, test.maskingCharacter, test.numberToMask)
			if got := buf.String(); got != test.want {
				t.Errorf("mask(%q, %s, %v) = %q, want %q", test.input, test.maskingCharacter, test.numberToMask, got, test.want)
			}
		})
	}
}

func TestDeidentifyDateShift(t *testing.T) {
	testutil.SystemTest(t)
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
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			deidentifyDateShift(buf, client, projectID, test.lowerBound, test.upperBound, test.input)
			if got := buf.String(); got != test.want {
				t.Errorf("deidentifyDateShift(%v, %v, %q) = %q, want %q", test.lowerBound, test.upperBound, got, test.input, test.want)
			}
		})
	}
}
