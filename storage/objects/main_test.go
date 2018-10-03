// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(m *testing.M) {
	// These functions are noisy.
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	os.Exit(s)
}
func TestObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var (
		bucket    = tc.ProjectID + "-samples-object-bucket-1"
		dstBucket = tc.ProjectID + "-samples-object-bucket-2"

		object1 = "foo.txt"
		object2 = "foo/a.txt"
	)

	cleanBucket(t, ctx, client, tc.ProjectID, bucket)
	cleanBucket(t, ctx, client, tc.ProjectID, dstBucket)

	if err := write(client, bucket, object1); err != nil {
		t.Fatalf("write(%q): %v", object1, err)
	}
	if err := write(client, bucket, object2); err != nil {
		t.Fatalf("write(%q): %v", object2, err)
	}

	{
		// Should only show "foo/a.txt", not "foo.txt"
		var buf bytes.Buffer
		if err := list(&buf, client, bucket); err != nil {
			t.Fatalf("cannot list objects: %v", err)
		}
		if got, want := buf.String(), object1; !strings.Contains(got, want) {
			t.Errorf("List() got %q; want to contain %q", got, want)
		}
		if got, want := buf.String(), object2; !strings.Contains(got, want) {
			t.Errorf("List() got %q; want to contain %q", got, want)
		}
	}

	{
		// Should only show "foo/a.txt", not "foo.txt"
		const prefix = "foo/"
		var buf bytes.Buffer
		if err := listByPrefix(&buf, client, bucket, prefix, ""); err != nil {
			t.Fatalf("cannot list objects by prefix: %v", err)
		}
		if got, want := buf.String(), object1; strings.Contains(got, want) {
			t.Errorf("List(%q) got %q; want NOT to contain %q", prefix, got, want)
		}
		if got, want := buf.String(), object2; !strings.Contains(got, want) {
			t.Errorf("List(%q) got %q; want to contain %q", prefix, got, want)
		}
	}

	data, err := read(client, bucket, object1)
	if err != nil {
		t.Fatalf("cannot read object: %v", err)
	}
	if got, want := string(data), "Hello\nworld"; got != want {
		t.Errorf("contents = %q; want %q", got, want)
	}
	_, err = attrs(client, bucket, object1)
	if err != nil {
		t.Errorf("cannot get object metadata: %v", err)
	}
	if err := makePublic(client, bucket, object1); err != nil {
		t.Errorf("cannot to make object public: %v", err)
	}
	err = move(client, bucket, object1)
	if err != nil {
		t.Fatalf("cannot move object: %v", err)
	}
	// object1's new name.
	object1 = object1 + "-rename"

	if err := copyToBucket(client, dstBucket, bucket, object1); err != nil {
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
	if err := addObjectACL(client, bucket, object1); err != nil {
		t.Errorf("cannot add object acl: %v", err)
	}
	if err := objectACL(client, bucket, object1); err != nil {
		t.Errorf("cannot get object acl: %v", err)
	}
	if err := objectACLFiltered(client, bucket, object1, storage.AllAuthenticatedUsers); err != nil {
		t.Errorf("cannot filter object acl: %v", err)
	}
	if err := deleteObjectACL(client, bucket, object1); err != nil {
		t.Errorf("cannot delete object acl: %v", err)
	}

	key := []byte("my-secret-AES-256-encryption-key")
	newKey := []byte("My-secret-AES-256-encryption-key")

	if err := writeEncryptedObject(client, bucket, object1, key); err != nil {
		t.Errorf("cannot write an encrypted object: %v", err)
	}
	data, err = readEncryptedObject(client, bucket, object1, key)
	if err != nil {
		t.Errorf("cannot read the encrypted object: %v", err)
	}
	if got, want := string(data), "top secret"; got != want {
		t.Errorf("object content = %q; want %q", got, want)
	}
	if err := rotateEncryptionKey(client, bucket, object1, key, newKey); err != nil {
		t.Errorf("cannot encrypt the object with the new key: %v", err)
	}
	if err := delete(client, bucket, object1); err != nil {
		t.Errorf("cannot to delete object: %v", err)
	}
	if err := delete(client, bucket, object2); err != nil {
		t.Errorf("cannot to delete object: %v", err)
	}

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		// Cleanup, this part won't be executed if Fatal happens.
		// TODO(jbd): Implement garbage cleaning.
		if err := client.Bucket(bucket).Delete(ctx); err != nil {
			r.Errorf("cleanup of bucket failed: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := delete(client, dstBucket, object1+"-copy"); err != nil {
			r.Errorf("cannot to delete copy object: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := client.Bucket(dstBucket).Delete(ctx); err != nil {
			r.Errorf("cleanup of bucket failed: %v", err)
		}
	})
}

func TestKMSObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	keyRingID := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKeyID := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")
	if keyRingID == "" || cryptoKeyID == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	var (
		bucket    = tc.ProjectID + "-samples-object-bucket-1"
		dstBucket = tc.ProjectID + "-samples-object-bucket-2"

		object1 = "foo.txt"
	)

	cleanBucket(t, ctx, client, tc.ProjectID, bucket)
	cleanBucket(t, ctx, client, tc.ProjectID, dstBucket)

	kmsKeyName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", tc.ProjectID, "global", keyRingID, cryptoKeyID)

	if err := writeWithKMSKey(client, bucket, object1, kmsKeyName); err != nil {
		t.Errorf("cannot write a KMS encrypted object: %v", err)
	}
}

func TestObjectBucketLock(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var (
		bucketName = tc.ProjectID + "-retent-samples-object-bucket"

		objectName = "foo.txt"

		retentionPeriod = 5 * time.Second
	)

	cleanBucket(t, ctx, client, tc.ProjectID, bucketName)
	bucket := client.Bucket(bucketName)

	if err := write(client, bucketName, objectName); err != nil {
		t.Fatalf("write(%q): %v", objectName, err)
	}
	if _, err := bucket.Update(ctx, storage.BucketAttrsToUpdate{
		RetentionPolicy: &storage.RetentionPolicy{
			RetentionPeriod: retentionPeriod,
		},
	}); err != nil {
		t.Errorf("unable to set retention policy (%q): %v", bucketName, err)
	}
	if err := setEventBasedHold(client, bucketName, objectName); err != nil {
		t.Errorf("unable to set event-based hold (%q/%q): %v", bucketName, objectName, err)
	}
	oAttrs, err := attrs(client, bucketName, objectName)
	if err != nil {
		t.Errorf("cannot get object metadata: %v", err)
	}
	if !oAttrs.EventBasedHold {
		t.Errorf("event-based hold is not enabled")
	}
	if err := releaseEventBasedHold(client, bucketName, objectName); err != nil {
		t.Errorf("unable to set event-based hold (%q/%q): %v", bucketName, objectName, err)
	}
	oAttrs, err = attrs(client, bucketName, objectName)
	if err != nil {
		t.Errorf("cannot get object metadata: %v", err)
	}
	if oAttrs.EventBasedHold {
		t.Errorf("event-based hold is not disabled")
	}
	if _, err := bucket.Update(ctx, storage.BucketAttrsToUpdate{
		RetentionPolicy: &storage.RetentionPolicy{},
	}); err != nil {
		t.Errorf("unable to remove retention policy (%q): %v", bucketName, err)
	}
	if err := setTemporaryHold(client, bucketName, objectName); err != nil {
		t.Errorf("unable to set temporary hold (%q/%q): %v", bucketName, objectName, err)
	}
	oAttrs, err = attrs(client, bucketName, objectName)
	if err != nil {
		t.Errorf("cannot get object metadata: %v", err)
	}
	if !oAttrs.TemporaryHold {
		t.Errorf("temporary hold is not disabled")
	}
	if err := releaseTemporaryHold(client, bucketName, objectName); err != nil {
		t.Errorf("unable to release temporary hold (%q/%q): %v", bucketName, objectName, err)
	}
	oAttrs, err = attrs(client, bucketName, objectName)
	if err != nil {
		t.Errorf("cannot get object metadata: %v", err)
	}
	if oAttrs.TemporaryHold {
		t.Errorf("temporary hold is not disabled")
	}
}

// cleanBucket ensures there's a fresh bucket with a given name, deleting the existing bucket if it already exists.
func cleanBucket(t *testing.T, ctx context.Context, client *storage.Client, projectID, bucket string) {
	b := client.Bucket(bucket)
	_, err := b.Attrs(ctx)
	if err == nil {
		it := b.Objects(ctx, nil)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("Bucket.Objects(%q): %v", bucket, err)
			}
			if attrs.EventBasedHold || attrs.TemporaryHold {
				if _, err := b.Object(attrs.Name).Update(ctx, storage.ObjectAttrsToUpdate{
					TemporaryHold:  false,
					EventBasedHold: false,
				}); err != nil {
					t.Fatalf("Bucket(%q).Object(%q).Update: %v", bucket, attrs.Name, err)
				}
			}
			if err := b.Object(attrs.Name).Delete(ctx); err != nil {
				t.Fatalf("Bucket(%q).Object(%q).Delete: %v", bucket, attrs.Name, err)
			}
		}
		if err := b.Delete(ctx); err != nil {
			t.Fatalf("Bucket.Delete(%q): %v", bucket, err)
		}
	}
	if err := b.Create(ctx, projectID, nil); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucket, err)
	}
}
