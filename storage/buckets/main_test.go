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
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"

	"cloud.google.com/go/storage"
)

var (
	storageClient *storage.Client
	bucketName    string
)

func TestMain(m *testing.M) {
	// These functions are noisy.
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	os.Exit(s)
}

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

func TestGetBucketMetadata(t *testing.T) {
	testutil.SystemTest(t)
	setup(t)

	bucketMetadataBuf := new(bytes.Buffer)
	if _, err := getBucketMetadata(bucketMetadataBuf, storageClient, bucketName); err != nil {
		t.Errorf("getBucketMetadata: %#v", err)
	}

	got := bucketMetadataBuf.String()
	if want := "BucketName:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
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

func TestBucketLock(t *testing.T) {
	tc := testutil.SystemTest(t)
	setup(t)

	retentionPeriod := 5 * time.Second

	if err := setRetentionPolicy(storageClient, bucketName, retentionPeriod); err != nil {
		t.Fatalf("failed to set retention policy (%q): %v", bucketName, err)
	}
	attrs, err := getRetentionPolicy(storageClient, bucketName)
	if err != nil {
		t.Fatalf("failed to get retention policy (%q): %v", bucketName, err)
	}
	if attrs.RetentionPolicy.RetentionPeriod != retentionPeriod {
		t.Fatalf("retention period is not the expected value (%q): %v", retentionPeriod, attrs.RetentionPolicy.RetentionPeriod)
	}
	if err := enableDefaultEventBasedHold(storageClient, bucketName); err != nil {
		t.Fatalf("failed to enable default event-based hold (%q): %v", bucketName, err)
	}
	attrs, err = getDefaultEventBasedHold(storageClient, bucketName)
	if err != nil {
		t.Fatalf("failed to get default event-based hold (%q): %v", bucketName, err)
	}
	if !attrs.DefaultEventBasedHold {
		t.Fatalf("default event-based hold was not enabled")
	}
	if err := disableDefaultEventBasedHold(storageClient, bucketName); err != nil {
		t.Fatalf("failed to disable event-based hold (%q): %v", bucketName, err)
	}
	attrs, err = getDefaultEventBasedHold(storageClient, bucketName)
	if err != nil {
		t.Fatalf("failed to get default event-based hold (%q): %v", bucketName, err)
	}
	if attrs.DefaultEventBasedHold {
		t.Fatalf("default event-based hold was not disabled")
	}
	if err := removeRetentionPolicy(storageClient, bucketName); err != nil {
		t.Fatalf("failed to remove retention policy (%q): %v", bucketName, err)
	}
	attrs, err = getRetentionPolicy(storageClient, bucketName)
	if err != nil {
		t.Fatalf("failed to get retention policy (%q): %v", bucketName, err)
	}
	if attrs.RetentionPolicy != nil {
		t.Fatalf("retention period to not be set")
	}
	if err := setRetentionPolicy(storageClient, bucketName, retentionPeriod); err != nil {
		t.Fatalf("failed to set retention policy (%q): %v", bucketName, err)
	}
	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := lockRetentionPolicy(storageClient, bucketName); err != nil {
			r.Errorf("failed to lock retention policy (%q): %v", bucketName, err)
		}
		attrs, err := getRetentionPolicy(storageClient, bucketName)
		if err != nil {
			r.Errorf("failed to check if retention policy is locked (%q): %v", bucketName, err)
		}
		if !attrs.RetentionPolicy.IsLocked {
			r.Errorf("retention policy is not locked")
		}
	})
	time.Sleep(5 * time.Second)
	deleteBucket(storageClient, bucketName)
	time.Sleep(5 * time.Second)
	if err := create(storageClient, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("failed to create bucket (%q): %v", bucketName, err)
	}
}

func TestUniformBucketLevelAccess(t *testing.T) {
	setup(t)

	if err := enableUniformBucketLevelAccess(storageClient, bucketName); err != nil {
		t.Fatalf("failed to enable uniform bucket-level access (%q): %v", bucketName, err)
	}
	attrs, err := getUniformBucketLevelAccess(storageClient, bucketName)
	if err != nil {
		t.Fatalf("failed to get uniform bucket-level access attrs (%q): %v", bucketName, err)
	}
	if !attrs.UniformBucketLevelAccess.Enabled {
		t.Fatalf("Uniform bucket-level access was not enabled for (%q).", bucketName)
	}
	if err := disableUniformBucketLevelAccess(storageClient, bucketName); err != nil {
		t.Fatalf("failed to disable uniform bucket-level access (%q): %v", bucketName, err)
	}
	attrs, err = getUniformBucketLevelAccess(storageClient, bucketName)
	if err != nil {
		t.Fatalf("failed to get uniform bucket-level access attrs (%q): %v", bucketName, err)
	}
	if attrs.UniformBucketLevelAccess.Enabled {
		t.Fatalf("Uniform bucket-level access was not disabled for (%q).", bucketName)
	}
}

func TestDelete(t *testing.T) {
	testutil.SystemTest(t)
	setup(t)

	if err := deleteBucket(storageClient, bucketName); err != nil {
		t.Fatalf("failed to delete bucket (%q): %v", bucketName, err)
	}
}
