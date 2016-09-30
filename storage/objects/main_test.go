// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	object := fmt.Sprintf("golang-example-objects-object-%d", time.Now().Unix())

	// TODO(jbd): Clean garbage buckets that are older than 1 day.
	bucket := fmt.Sprintf("golang-example-objects-bucket-%d", time.Now().Unix())
	dstBucket := fmt.Sprintf("golang-example-objects-dstbucket-%d", time.Now().Unix())

	ensureBucketExists(ctx, client, tc.ProjectID, bucket)
	ensureBucketExists(ctx, client, tc.ProjectID, dstBucket)

	if err := write(client, bucket, object); err != nil {
		t.Fatalf("cannot write object: %v", err)
	}
	data, err := read(client, bucket, object)
	if err != nil {
		t.Fatalf("cannot read object: %v", err)
	}
	if got, want := string(data), "Hello\nworld"; got != want {
		t.Errorf("contents = %q; want %q", got, want)
	}
	_, err = attrs(client, bucket, object)
	if err != nil {
		t.Errorf("cannot get object metadata: %v", err)
	}
	if err := makePublic(client, bucket, object); err != nil {
		t.Errorf("cannot to make object public: %v", err)
	}
	err = move(client, bucket, object)
	if err != nil {
		t.Fatalf("cannot move object: %v", err)
	}
	object += "-rename"
	err = copyToBucket(client, dstBucket, bucket, object)
	if err != nil {
		t.Errorf("cannot copy object to bucket: %v", err)
	}
	if err := addBucketACL(client, bucket); err != nil {
		t.Errorf("cannot add bucket acl: %v", err)
	}
	if err := addDefaultBucketACL(client, bucket); err != nil {
		t.Errorf("cannot add bucket deafult acl: %v", err)
	}
	if err := bucketACL(client, bucket); err != nil {
		t.Errorf("cannot get bucket acl: %v", err)
	}
	if err := bucketACLFiltered(client, bucket, storage.AllAuthenticatedUsers); err != nil {
		t.Errorf("cannot filter bucket acl: %v", err)
	}
	if err := deleteDefaultBucketACL(client, bucket); err != nil {
		t.Errorf("cannot delete bucket default acl: %v", err)
	}
	if err := deleteBucketACL(client, bucket); err != nil {
		t.Errorf("cannot delete bucket acl: %v", err)
	}
	if err := addObjectACL(client, bucket, object); err != nil {
		t.Errorf("cannot add object acl: %v", err)
	}
	if err := objectACL(client, bucket, object); err != nil {
		t.Errorf("cannot get object acl: %v", err)
	}
	if err := objectACLFiltered(client, bucket, object, storage.AllAuthenticatedUsers); err != nil {
		t.Errorf("cannot filter object acl: %v", err)
	}
	if err := deleteObjectACL(client, bucket, object); err != nil {
		t.Errorf("cannot delete object acl: %v", err)
	}
	if err := delete(client, bucket, object); err != nil {
		t.Errorf("cannot to delete object: %v", err)
	}

	// Cleanup, this part won't be executed if Fatal happens.
	// TODO(jbd): Implement garbage cleaning.
	if err := client.Bucket(bucket).Delete(ctx); err != nil {
		t.Fatalf("cleanup of bucket failed: %v", err)
	}

	time.Sleep(4 * time.Second) // For eventual consistency.
	if err := delete(client, dstBucket, object+"-copy"); err != nil {
		t.Errorf("cannot to delete copy object: %v", err)
	}
	if err := client.Bucket(dstBucket).Delete(ctx); err != nil {
		t.Fatalf("cleanup of bucket failed: %v", err)
	}
}

func ensureBucketExists(ctx context.Context, client *storage.Client, projectID, bucket string) {
	b := client.Bucket(bucket)
	_, err := b.Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		err = b.Create(ctx, projectID, nil)
	}
	if err != nil {
		log.Fatalf("bucket ensuring failed: %v", err)
	}
}
