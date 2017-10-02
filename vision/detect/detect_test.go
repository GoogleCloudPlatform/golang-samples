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
		{"Labels2", detectLabels, detectLabelsURI, "wakeupcat.jpg", "whiskers"},
		{"Faces", detectFaces, detectFacesURI, "face_no_surprise.jpg", "Anger"},
		{"Landmarks", detectLandmarks, detectLandmarksURI, "landmark.jpg", "Palace"},
		{"Logos", detectLogos, detectLogosURI, "logos.png", "Google"},
		{"Properties", detectProperties, detectPropertiesURI, "landmark.jpg", "%"},
		{"SafeSearch", detectSafeSearch, detectSafeSearchURI, "wakeupcat.jpg", "Spoofed"},
		{"Text", detectText, detectTextURI, "text.jpg", "Preparing to install"},
		{"FullText", detectDocumentText, detectDocumentTextURI, "text.jpg", "Preparing to install"},
		{"Crop", detectCropHints, detectCropHintsURI, "wakeupcat.jpg", "(0,0)"},
		{"Web", detectWeb, detectWebURI, "wakeupcat.jpg", "Web properties"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		err := tt.local(&buf, "../testdata/"+tt.path)
		if err != nil {
			t.Fatalf("Local %s(%q): got %v, want nil err", tt.name, tt.path, err)
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
			t.Fatalf("GCS %s(%q): got %v, want nil err", tt.name, tt.path, err)
		}
		if got := buf.String(); !strings.Contains(got, tt.wantContain) {
			t.Errorf("GCS %s(%q): got %q, want to contain %q", tt.name, tt.path, got, tt.wantContain)
		}
	}
}
