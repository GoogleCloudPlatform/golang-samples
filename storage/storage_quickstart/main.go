// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START storage_quickstart]
// Sample storage_quickstart creates a Google Cloud Storage bucket.
package main

import (
	"fmt"
	"golang.org/x/net/context"
	"log"

	// Imports the Google Cloud Storage client package
	"cloud.google.com/go/storage"
)

func main() {
	ctx := context.Background()

	// Your Google Cloud Platform project ID
	projectID := "YOUR_PROJECT_ID"

	// Creates a client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// The name for the new bucket
	bucketName := "my-new-bucket"

	// Prepares a new bucket
	bucket := client.Bucket(bucketName)

	// Creates the new bucket
	if err := bucket.Create(ctx, projectID, nil); err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	fmt.Printf("Bucket %v created.", bucketName)
}

// [END storage_quickstart]
