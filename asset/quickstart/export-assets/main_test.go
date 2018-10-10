// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := fmt.Sprintf("%s-for-assets", tc.ProjectID)
	os.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create storage client: %v", err)
	}

	// Delete the bucket (if it exists) then recreate it.
	cleanBucket(ctx, t, client, tc.ProjectID, bucketName)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

	want := fmt.Sprintf("output_config:<gcs_destination:<uri:\"gs://%s/my-assets.txt\" > >", bucketName)
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
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
