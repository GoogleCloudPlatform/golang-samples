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

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
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
		{"Objects", localizeObjects, localizeObjectsURI, "puppies.jpg", "Dog"},
	}

	for _, tt := range tests {
		if tt.local == nil {
			continue
		}
		t.Run(tt.name+"/local", func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			if err := tt.local(&buf, "../testdata/"+tt.path); err != nil {
				t.Fatalf("Local %s(%q): got %v, want nil err", tt.name, tt.path, err)
			}
			if got, wantContain := strings.ToLower(buf.String()), strings.ToLower(tt.wantContain); !strings.Contains(got, wantContain) {
				t.Errorf("Local %s(%q): got %q, want to contain %q", tt.name, tt.path, got, wantContain)
			}
		})
	}

	for _, tt := range tests {
		if tt.gcs == nil {
			continue
		}
		t.Run(tt.name+"/gcs", func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			if err := tt.gcs(&buf, "gs://python-docs-samples-tests/vision/"+tt.path); err != nil {
				t.Fatalf("GCS %s(%q): got %v, want nil err", tt.name, tt.path, err)
			}
			if got, wantContain := strings.ToLower(buf.String()), strings.ToLower(tt.wantContain); !strings.Contains(got, wantContain) {
				t.Errorf("GCS %s(%q): got %q, want to contain %q", tt.name, tt.path, got, wantContain)
			}
		})
	}
}

func TestDetectAsyncDocument(t *testing.T) {
	tc := testutil.SystemTest(t)

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	bucketName := fmt.Sprintf("%s-vision", tc.ProjectID)
	bucket := client.Bucket(bucketName)
	cleanBucket(ctx, t, client, tc.ProjectID, bucketName)

	var buf bytes.Buffer
	gcsSourceURI := "gs://python-docs-samples-tests/HodgeConj.pdf"
	gcsDestinationURI := "gs://" + bucketName + "/vision/"
	err = detectAsyncDocumentURI(&buf, gcsSourceURI, gcsDestinationURI)
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

func cleanBucket(ctx context.Context, t *testing.T, client *storage.Client, projectID, bucket string) {
	deleteBucketIfExists(ctx, t, client, bucket)

	b := client.Bucket(bucket)
	// Now create it
	if err := b.Create(ctx, projectID, nil); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucket, err)
	}
}

func deleteBucketIfExists(ctx context.Context, t *testing.T, client *storage.Client, bucket string) {
	b := client.Bucket(bucket)
	if _, err := b.Attrs(ctx); err != nil {
		return
	}

	// Delete all the elements in the already existent bucket
	it := b.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("Bucket.Objects(%q): %v", bucket, err)
		}
		if err := b.Object(attrs.Name).Delete(ctx); err != nil {
			t.Fatalf("Bucket(%q).Object(%q).Delete: %v", bucket, attrs.Name, err)
		}
	}
	// Then delete the bucket itself
	if err := b.Delete(ctx); err != nil {
		t.Fatalf("Bucket.Delete(%q): %v", bucket, err)
	}
}
