// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"

	"cloud.google.com/go/storage"
)

var bucketName string

var setupOnce sync.Once

func setup(t *testing.T) *storage.Client {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	setupOnce.Do(func() {
		bucketName = fmt.Sprintf("golang-example-bucketmgr-%d", time.Now().Unix())
	})
	return client
}

func TestCreate(t *testing.T) {
	tc := testutil.SystemTest(t)
	c := setup(t)
	if err := create(c, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("failed to create bucket (%q): %v", bucketName, err)
	}
	time.Sleep(3 * time.Second) // for eventual consistency
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)
	c := setup(t)
	buckets, err := list(c, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	var ok bool
	for _, b := range buckets {
		if b == bucketName {
			ok = true
			break
		}
	}
	if !ok {
		t.Errorf("got bucket list: %+v; want %q in the list", strings.Join(buckets, "\n"), bucketName)
	}
}

func TestDelete(t *testing.T) {
	c := setup(t)
	if err := delete(c, bucketName); err != nil {
		t.Fatalf("failed to delete bucket (%q): %v", bucketName, err)
	}
}
