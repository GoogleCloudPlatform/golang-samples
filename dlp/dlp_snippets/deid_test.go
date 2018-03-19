// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"
)

func TestMask(t *testing.T) {
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
		buf := new(bytes.Buffer)
		mask(buf, client, projectID, test.input, test.maskingCharacter, test.numberToMask)
		if got := buf.String(); got != test.want {
			t.Errorf("mask(%q, %s, %v) = %q, want %q", test.input, test.maskingCharacter, test.numberToMask, got, test.want)
		}
	}
}

func TestDeidentifyDateShift(t *testing.T) {
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
		buf := new(bytes.Buffer)
		deidentifyDateShift(buf, client, projectID, test.lowerBound, test.upperBound, test.input)
		if got := buf.String(); got != test.want {
			t.Errorf("deidentifyDateShift(%v, %v, %q) = %q, want %q", test.lowerBound, test.upperBound, got, test.input, test.want)
		}
	}
}
