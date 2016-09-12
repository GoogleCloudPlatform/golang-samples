// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command bucketmgr is a program that manages Google Storage buckets by using
// the Google Storage API. More documentation is available at
// https://cloud.google.com/storage/docs/json_api/v1/.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func main() {
	ctx := context.Background()

	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}

	// [START setup]
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// [END setup]

	// Give the bucket a unique name.
	name := fmt.Sprintf("golang-example-bucketmgr-%d", time.Now().Unix())
	if err := create(client, proj, name); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created bucket: %v\n", name)

	// list buckets from the project
	buckets, err := list(client, proj)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("buckets: %+v\n", buckets)

	// delete the bucket
	if err := delete(client, name); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted bucket: %v\n", name)
}

func create(client *storage.Client, proj, name string) error {
	ctx := context.Background()
	// [START create_bucket]
	if err := client.Bucket(name).Create(ctx, proj, nil); err != nil {
		return err
	}
	// [END create_bucket]
	// Do not merge with the previous if, we are trying to make the
	// tagged snippets look good.
	return nil
}

func list(client *storage.Client, proj string) ([]string, error) {
	ctx := context.Background()
	// [START list_buckets]
	var buckets []string
	it := client.Buckets(ctx, proj)
	for {
		battrs, err := it.Next()
		if err == storage.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, battrs.Name)
	}
	// [END list_buckets]
	return buckets, nil
}

func delete(client *storage.Client, name string) error {
	ctx := context.Background()
	// [START delete_bucket]
	if err := client.Bucket(name).Delete(ctx); err != nil {
		return err
	}
	// [END delete_bucket]
	// Do not merge with the previous if, we are trying to make the
	// tagged snippets look good.
	return nil
}
