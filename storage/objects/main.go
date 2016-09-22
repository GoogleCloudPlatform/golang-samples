// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample objects creates, list, deletes objects and runs
// other similar operations on them by using the Google Storage API.
// More documentation is available at
// https://cloud.google.com/storage/docs/json_api/v1/.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"cloud.google.com/go/storage"
)

var objectName = fmt.Sprintf("golang-example-objects-%d", time.Now().Unix())

const (
	bucketName    = "golang-example-objects-bucket"
	dstBucketName = "golang-example-objects-bucket-dst"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := createBucketIfNotexits(ctx, client, projectID, bucketName); err != nil {
		log.Fatal(err)
	}
	// create another bucket later to be used to copy objects from the prev one.
	if err := createBucketIfNotexits(ctx, client, projectID, dstBucketName); err != nil {
		log.Fatal(err)
	}
	name := fmt.Sprintf("gs://%v/%v", bucketName, objectName)

	if err := write(client, bucketName, objectName); err != nil {
		log.Fatalf("failed to write to the object (%v): %v", objectName, err)
	}

	data, err := read(client, bucketName, objectName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("contents of %v: %q\n", name, string(data))

	attrs, err := attrs(client, bucketName, objectName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("object attributes of %v: %+v\n", name, attrs)

	err = makePublic(client, bucketName, objectName)
	if err != nil {
		log.Fatalf("failed to make the object public: %v", err)
	}

	newName, err := move(client, bucketName, objectName)
	if err != nil {
		log.Fatalf("failed to move the object: %v", err)
	}
	objectName = newName

	if err := delete(client, bucketName, objectName); err != nil {
		log.Fatalf("failed to delete the object: %v", err)
	}
}

func write(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	// [START upload_file]
	f, err := os.Open("notes.txt")
	if err != nil {
		return err
	}
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END upload_file]
	return nil
}

func read(client *storage.Client, bucket, object string) ([]byte, error) {
	ctx := context.Background()
	// [START download_file]
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	if err := rc.Close(); err != nil {
		return nil, err
	}
	return data, nil
	// [END download_file]
}

func attrs(client *storage.Client, bucket, object string) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	// [START get_metadata]
	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	return attrs, nil
	// [END get_metadata]
}

func makePublic(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	// [START public]
	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}
	// [END public]
	return nil
}

func move(client *storage.Client, bucket, object string) (string, error) {
	ctx := context.Background()
	// [START move_file]
	dstName := object + "-rename"

	src := client.Bucket(bucket).Object(object)
	dst := client.Bucket(bucket).Object(dstName)

	if _, err := src.CopyTo(ctx, dst, nil); err != nil {
		return "", err
	}
	if err := src.Delete(ctx); err != nil {
		return "", err
	}
	// [END move_file]
	return dstName, nil
}

func copyToBucket(client *storage.Client, dstBucket, srcBucket, srcObject string) error {
	ctx := context.Background()
	// [START copy_file]
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(srcObject + "-copy")

	if _, err := src.CopyTo(ctx, dst, nil); err != nil {
		return err
	}
	// [END copy_file]
	return nil
}

func delete(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	// [START download_file]
	o := client.Bucket(bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	// [END download_file]
	return nil
}

func createBucketIfNotexits(ctx context.Context, client *storage.Client, projectID, bucket string) error {
	b := client.Bucket(bucket)
	_, err := b.Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		err = b.Create(ctx, projectID, nil)
	}
	return err
}
