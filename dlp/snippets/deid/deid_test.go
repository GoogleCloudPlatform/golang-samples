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
	"strings"
	"testing"

	"cloud.google.com/go/dlp/apiv2/dlppb"
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

func TestDeidentifyTableConditionInfoTypes(t *testing.T) {
	tc := testutil.SystemTest(t)

	row1 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "22"}},
			{Type: &dlppb.Value_StringValue{StringValue: "Jane Austen"}},
			{Type: &dlppb.Value_StringValue{StringValue: "21"}},
			{Type: &dlppb.Value_StringValue{StringValue: "There are 14 kisses in Jane Austen's novels."}},
		},
	}

	row2 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "55"}},
			{Type: &dlppb.Value_StringValue{StringValue: "Mark Twain"}},
			{Type: &dlppb.Value_StringValue{StringValue: "75"}},
			{Type: &dlppb.Value_StringValue{StringValue: "Mark Twain loved cats."}},
		},
	}

	row3 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "101"}},
			{Type: &dlppb.Value_StringValue{StringValue: "Charles Dickens"}},
			{Type: &dlppb.Value_StringValue{StringValue: "95"}},
			{Type: &dlppb.Value_StringValue{StringValue: "Charles Dickens name was a curse invented by Shakespeare."}},
		},
	}

	table := &dlppb.Table{
		Headers: []*dlppb.FieldId{
			{Name: "AGE"},
			{Name: "PATIENT"},
			{Name: "HAPPINESS SCORE"},
			{Name: "FACTOID"},
		},
		Rows: []*dlppb.Table_Row{
			{Values: row1.Values},
			{Values: row2.Values},
			{Values: row3.Values},
		},
	}

	buf := new(bytes.Buffer)

	if err := deidentifyTableConditionInfoTypes(buf, tc.ProjectID, table, []string{"PATIENT", "FACTOID"}); err != nil {
		t.Errorf("deidentifyTableConditionInfoTypes: %v", err)
	}

	got := buf.String()
	if want := "Table after de-identification"; !strings.Contains(got, want) {
		t.Errorf("deidentifyTableConditionInfoTypes got %q, want %q", got, want)
	}
}
