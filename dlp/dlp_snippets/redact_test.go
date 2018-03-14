/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"strings"
	"testing"

	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func TestRedactImage(t *testing.T) {
	tests := []struct {
		name      string
		inputPath string
		bt        dlppb.ByteContentItem_BytesType
		infoTypes []string
		want      string
	}{
		{
			name:      "image with one type",
			inputPath: "testdata/ok.png",
			bt:        dlppb.ByteContentItem_IMAGE_PNG,
			infoTypes: []string{"US_SOCIAL_SECURITY_NUMBER"},
			want:      "Wrote output to",
		},
		{
			name:      "image with two types",
			inputPath: "testdata/ok.png",
			bt:        dlppb.ByteContentItem_IMAGE_PNG,
			infoTypes: []string{"US_SOCIAL_SECURITY_NUMBER", "DATE"},
			want:      "Wrote output to",
		},
	}
	for _, test := range tests {
		buf := new(bytes.Buffer)
		redactImage(buf, client, projectID, dlppb.Likelihood_POSSIBLE, test.infoTypes, test.bt, test.inputPath, "testdata/test_output.png")
		if got := buf.String(); !strings.Contains(got, test.want) {
			t.Errorf("redactImage(%s) got %q, want substring %q", test.name, got, test.want)
		}
	}
}
