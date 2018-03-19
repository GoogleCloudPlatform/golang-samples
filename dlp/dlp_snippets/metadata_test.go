// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestInfoTypes(t *testing.T) {
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
		buf := new(bytes.Buffer)
		infoTypes(buf, client, test.language, test.filter)
		if got := buf.String(); !strings.Contains(got, test.want) {
			t.Errorf("infoTypes(%s, %s) = %s, want substring %q", test.language, test.filter, got, test.want)
		}
	}
}
