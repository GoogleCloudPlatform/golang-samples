// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

var (
	ns     = testutil.Namespace("buckets")
	bucket = testutil.NextBucket(ns)
)

func setup(t *testing.T) *storage.Client {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func TestCreate(t *testing.T) {
	tc := testutil.SystemTest(t)
	c := setup(t)
	if err := create(c, tc.ProjectID, bucket); err != nil {
		t.Fatalf("failed to create bucket (%q): %v", bucket, err)
	}
}

func TestCreateWithAttrs(t *testing.T) {
	tc := testutil.SystemTest(t)
	c := setup(t)
	name := testutil.NextBucket("buckets")
	if err := createWithAttrs(c, tc.ProjectID, name); err != nil {
		t.Fatalf("failed to create bucket (%q): %v", bucket, err)
	}
	if err := delete(c, name); err != nil {
		t.Fatalf("failed to delete bucket (%q): %v", bucket, err)
	}
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)
	c := setup(t)
	buckets, err := list(c, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	var ok bool
outer:
	for attempt := 0; attempt < 5; attempt++ { // for eventual consistency
		for _, b := range buckets {
			if b == bucket {
				ok = true
				break outer
			}
		}
		time.Sleep(2 * time.Second)
	}
	if !ok {
		t.Errorf("got bucket list: %v; want %q in the list", buckets, bucket)
	}
}

func TestDelete(t *testing.T) {
	testutil.SystemTest(t)
	defer testutil.CleanupBuckets(t, bucket)

	c := setup(t)
	if err := delete(c, bucket); err != nil {
		t.Fatalf("failed to delete bucket (%q): %v", bucket, err)
	}
}
