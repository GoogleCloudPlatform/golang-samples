// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"

	"cloud.google.com/go/storage"
)

var (
	testBucket     = fmt.Sprintf("golang-example-objects-bucket-%d", time.Now().Unix())
	testDestBucket = fmt.Sprintf("golang-example-objects-bucket-dest-%d", time.Now().Unix())
	testObject     = fmt.Sprintf("golang-example-object-%d", time.Now().Unix())
)

var bucketOnce sync.Once

func setup(t *testing.T) *storage.Client {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	bucketOnce.Do(func() {
		// TODO(jbd): delete buckets older than a day
		if err := createBucketIfNotexits(ctx, client, tc.ProjectID, testBucket); err != nil {
			t.Fatal(err)
		}
		if err := createBucketIfNotexits(ctx, client, tc.ProjectID, testDestBucket); err != nil {
			t.Fatal(err)
		}
	})
	return client
}

func TestWrite(t *testing.T) {
	c := setup(t)

	if err := write(c, testBucket, testObject); err != nil {
		t.Errorf("failed to write: %v", err)
	}
}

func TestRead(t *testing.T) {
	c := setup(t)

	data, err := read(c, testBucket, testObject)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(data), "hello\nworld"; got != want {
		t.Errorf("object contents = %q; want %q", got, want)
	}
}

func TestAttrs(t *testing.T) {
	c := setup(t)

	_, err := attrs(c, testBucket, testObject)
	if err != nil {
		t.Errorf("failed to get object attributes: %v", err)
	}
}

func TestMakePublic(t *testing.T) {
	c := setup(t)

	if err := makePublic(c, testBucket, testObject); err != nil {
		t.Errorf("failed to make public: %v", err)
	}
}

func TestMove(t *testing.T) {
	c := setup(t)

	name, err := move(c, testBucket, testObject)
	if err != nil {
		t.Fatalf("failed to move: %v", err)
	}
	testObject = name
}

func TestCopyToBucket(t *testing.T) {
	c := setup(t)

	err := copyToBucket(c, testDestBucket, testBucket, testObject)
	if err != nil {
		t.Fatalf("failed to move: %v", err)
	}

	if err := delete(c, testDestBucket, testObject+"-copy"); err != nil {
		t.Errorf("failed to delete object (%q): %v", testObject, err)
	}
}

func TestDelete(t *testing.T) {
	c := setup(t)

	if err := delete(c, testBucket, testObject); err != nil {
		t.Errorf("failed to delete object (%q): %v", testObject, err)
	}
}
