// Copyright 2017 Google Inc. All rights reserved.
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

const catVideo = "gs://demomaker/cat.mp4"

func TestAnalyze(t *testing.T) {
	testutil.SystemTest(t)

	tests := []struct {
		name        string
		gcs         func(io.Writer, string) error
		path        string
		wantContain string
	}{
		{"ShotChange", shotChangeURI, catVideo, "Shot"},
		{"Labels", labelURI, catVideo, "cat"},
		{"Explicit", explicitContentURI, catVideo, "VERY_UNLIKELY"},
	}

	for _, tt := range tests {
		if tt.gcs == nil {
			continue
		}

		var buf bytes.Buffer
		err := tt.gcs(&buf, tt.path)
		if err != nil {
			t.Fatalf("GCS %s(%q): got %v, want nil err", tt.name, tt.path, err)
		}
		if got := buf.String(); !strings.Contains(got, tt.wantContain) {
			t.Errorf("GCS %s(%q): got %q, want to contain %q", tt.name, tt.path, got, tt.wantContain)
		}
	}
}

func TestGenerated(t *testing.T) {
	testutil.Generated(t, "gen/template.go").
		Goimports().
		Matches("video_analyze.go")

	testutil.Generated(t, "gen/template.go").
		Labels("gcs").
		Goimports().
		Matches("video_analyze_gcs.go")
}
