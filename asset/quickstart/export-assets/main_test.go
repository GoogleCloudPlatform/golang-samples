// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"

	"fmt"
	"io/ioutil"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

var (
	bucketName string
)

func setup(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	var storageClient *storage.Client
	var err error
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create storage client: %v", err)
	}

	existing, err := CheckBucketExistance(tc.ProjectID, storageClient)
	if err != nil {
		t.Fatalf("failed to list buckets: %v", err)
	}
	if !existing {
		if err := storageClient.Bucket(bucketName).Create(ctx, tc.ProjectID, nil); err != nil {
			t.Fatalf("failed to create bucket (%q): %v", bucketName, err)
		}
	}
}

func CheckBucketExistance(projectID string, storageClient *storage.Client) (bool, error) {
	ctx := context.Background()
	it := storageClient.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, err
		}
		if battrs.Name == bucketName {
			return true, nil
		}
	}
	return false, nil
}

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	os.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	bucketName = fmt.Sprintf("%s-for-assets", tc.ProjectID)

	setup(t)

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
