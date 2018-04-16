// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func TestRedactImage(t *testing.T) {
	testutil.SystemTest(t)
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
