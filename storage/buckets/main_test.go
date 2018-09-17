// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"

	"cloud.google.com/go/storage"
)

var (
	storageClient *storage.Client
	bucketName    string
)

func setup(t *testing.T) {
	tc := testutil.SystemTest(t)

	ctx := context.Background()
	var err error
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	bucketName = tc.ProjectID + "-storage-buckets-tests"
}

func TestCreate(t *testing.T) {
	tc := testutil.SystemTest(t)
	setup(t)

	// Clean up bucket before running tests.
	deleteBucket(storageClient, bucketName)
	if err := create(storageClient, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("failed to create bucket (%q): %v", bucketName, err)
	}
}

func TestCreateWithAttrs(t *testing.T) {
	tc := testutil.SystemTest(t)
	name := bucketName + "-attrs"
	// Clean up bucket before running test.
	deleteBucket(storageClient, name)
	if err := createWithAttrs(storageClient, tc.ProjectID, name); err != nil {
		t.Fatalf("failed to create bucket (%q): %v", name, err)
	}
	if err := deleteBucket(storageClient, name); err != nil {
		t.Fatalf("failed to delete bucket (%q): %v", name, err)
	}
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)
	setup(t)

	buckets, err := list(storageClient, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	var ok bool
outer:
	for attempt := 0; attempt < 5; attempt++ { // for eventual consistency
		for _, b := range buckets {
			if b == bucketName {
				ok = true
				break outer
			}
		}
		time.Sleep(2 * time.Second)
	}
	if !ok {
		t.Errorf("got bucket list: %v; want %q in the list", buckets, bucketName)
	}
}

func TestIAM(t *testing.T) {
	testutil.SystemTest(t)
	setup(t)

	if _, err := getPolicy(storageClient, bucketName); err != nil {
		t.Errorf("getPolicy: %#v", err)
	}
	if err := addUser(storageClient, bucketName); err != nil {
		t.Errorf("addUser: %v", err)
	}
	if err := removeUser(storageClient, bucketName); err != nil {
		t.Errorf("removeUser: %v", err)
	}
}

func TestRequesterPays(t *testing.T) {
	testutil.SystemTest(t)
	setup(t)

	if err := enableRequesterPays(storageClient, bucketName); err != nil {
		t.Errorf("enableRequesterPays: %#v", err)
	}
	if err := disableRequesterPays(storageClient, bucketName); err != nil {
		t.Errorf("disableRequesterPays: %#v", err)
	}
	if err := checkRequesterPays(storageClient, bucketName); err != nil {
		t.Errorf("checkRequesterPays: %#v", err)
	}
}

func TestKMS(t *testing.T) {
	tc := testutil.SystemTest(t)
	setup(t)

	keyRingID := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKeyID := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")

	if keyRingID == "" || cryptoKeyID == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	kmsKeyName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", tc.ProjectID, "global", keyRingID, cryptoKeyID)
	if err := setDefaultKMSkey(storageClient, bucketName, kmsKeyName); err != nil {
		t.Fatalf("failed to enable default kms key (%q): %v", bucketName, err)
	}
}

func TestDelete(t *testing.T) {
	testutil.SystemTest(t)
	setup(t)

	if err := deleteBucket(storageClient, bucketName); err != nil {
		t.Fatalf("failed to delete bucket (%q): %v", bucketName, err)
	}
}
