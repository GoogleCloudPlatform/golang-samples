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

// Package deid contains example snippets using the DLP deidentification API.
package deid

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMask(t *testing.T) {
	tc := testutil.SystemTest(t)
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
		test := test
		t.Run(test.input, func(t *testing.T) {
			test := test
			t.Parallel()
			buf := new(bytes.Buffer)
			err := mask(buf, tc.ProjectID, test.input, []string{"US_SOCIAL_SECURITY_NUMBER"}, test.maskingCharacter, test.numberToMask)
			if err != nil {
				t.Errorf("mask(%q, %s, %v) = error %q, want %q", test.input, test.maskingCharacter, test.numberToMask, err, test.want)
			}
			if got := buf.String(); got != test.want {
				t.Errorf("mask(%q, %s, %v) = %q, want %q", test.input, test.maskingCharacter, test.numberToMask, got, test.want)
			}
		})
	}
}

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

func TestDeIdentifyWithRedact(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "My name is Alicia Abernathy, and my email address is aabernathy@example.com."
	infoTypeNames := []string{"EMAIL_ADDRESS"}
	want := "output: My name is Alicia Abernathy, and my email address is ."

	buf := new(bytes.Buffer)
	err := deidentifyWithRedact(buf, tc.ProjectID, input, infoTypeNames)
	if err != nil {
		t.Errorf("deidentifyWithRedact(%q) = error '%q', want %q", err, input, want)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyWithRedact(%q) = %q, want %q", got, input, want)
	}
}

func TestDeidentifyExceptionList(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "jack@example.org accessed customer record of user5@example.com"
	want := "output : jack@example.org accessed customer record of [EMAIL_ADDRESS]"

	buf := new(bytes.Buffer)
	err := deidentifyExceptionList(buf, tc.ProjectID, input)
	if err != nil {
		t.Errorf("deidentifyExceptionList(%q) = error '%q', want %q", input, err, want)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyExceptionList(%q) = %q, want %q", input, got, want)
	}
}

func TestDeIdentifyWithWordList(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "Patient was seen in RM-YELLOW then transferred to rm green."
	infoType := "CUSTOM_ROOM_ID"
	wordList := []string{"RM-GREEN", "RM-YELLOW", "RM-ORANGE"}
	want := "output : Patient was seen in [CUSTOM_ROOM_ID] then transferred to [CUSTOM_ROOM_ID]."

	buf := new(bytes.Buffer)
	err := deidentifyWithWordList(buf, tc.ProjectID, input, infoType, wordList)
	if err != nil {
		t.Errorf("deidentifyWithWordList(%q) = error '%q', want %q", input, err, want)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyWithWordList(%q) = %q, want %q", input, got, want)
	}
}
