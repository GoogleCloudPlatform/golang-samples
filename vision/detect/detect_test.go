// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDetect(t *testing.T) {
	testutil.SystemTest(t)

	tests := []struct {
		name        string
		local, gcs  func(io.Writer, string) error
		path        string
		wantContain string
	}{
		{"Labels", detectLabels, nil, "cat.jpg", "cat"},
		{"Labels2", detectLabels, detectLabelsGCS, "wakeupcat.jpg", "whiskers"},
		{"Faces", detectFaces, detectFacesGCS, "face_no_surprise.jpg", "Anger"},
		{"Landmarks", detectLandmarks, detectLandmarksGCS, "landmark.jpg", "Palace"},
		{"Logos", detectLogos, detectLogosGCS, "logos.png", "Google"},
		{"Properties", detectProperties, detectPropertiesGCS, "landmark.jpg", "%"},
		{"SafeSearch", detectSafeSearch, detectSafeSearchGCS, "wakeupcat.jpg", "Spoofed"},
		{"Text", detectText, detectTextGCS, "text.jpg", "Preparing to install"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		err := tt.local(&buf, "../testdata/"+tt.path)
		if err != nil {
			t.Fatalf("got %v, want nil err", err)
		}
		if got := buf.String(); !strings.Contains(got, tt.wantContain) {
			t.Errorf("Local %s(%q): got %q, want to contain %q", tt.name, tt.path, got, tt.wantContain)
		}
	}

	for _, tt := range tests {
		if tt.gcs == nil {
			continue
		}

		var buf bytes.Buffer
		err := tt.gcs(&buf, "gs://python-docs-samples-tests/vision/"+tt.path)
		if err != nil {
			t.Fatalf("got %v, want nil err", err)
		}
		if got := buf.String(); !strings.Contains(got, tt.wantContain) {
			t.Errorf("GCS %s(%q): got %q, want to contain %q", tt.name, tt.path, got, tt.wantContain)
		}
	}
}
