// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"

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
		{"WebGeo", nil, detectWebGeoURI, "city.jpg", "Entities"},
	}

	for _, tt := range tests {
		if tt.local == nil {
			continue
		}

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

func TestDetectAsyncDocument(t *testing.T) {
	tc := testutil.SystemTest(t)

	ctx := context.Background()

	// Create a temporary bucket
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	bucketName := fmt.Sprintf("%s-golang-samples-%d", tc.ProjectID, time.Now().Unix())
	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, tc.ProjectID, nil); err != nil {
		t.Fatal(err)
	}

	// Clean and delete the bucket at the end of the test
	defer func() {
		it := bucket.Objects(ctx, nil)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
				t.Fatal(err)
			}
		}
		if err := bucket.Delete(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	// Run the test
	var buf bytes.Buffer
	gcsSourceURI := "gs://python-docs-samples-tests/HodgeConj.pdf"
	gcsDestinationURI := "gs://" + bucketName + "/vision/"
	err = detectAsyncDocument(&buf, gcsSourceURI, gcsDestinationURI)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output files exist
	expectedFiles := []string{
		"vision/output-1-to-2.json",
		"vision/output-3-to-4.json",
		"vision/output-5-to-5.json",
	}
	for _, filename := range expectedFiles {
		_, err = bucket.Object(filename).Attrs(ctx)
		if err != nil {
			t.Fatalf("wanted object %q, got error: %v", filename, err)
		}
	}
}
