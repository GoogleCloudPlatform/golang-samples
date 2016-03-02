// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
)
import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

// ListBuckets returns a slice of all the buckets in the given projectId.
// [START ListBuckets]
func ListBuckets(projectId string) ([]*storage.Bucket, error) {
	// Create the client that uses Application Default Credentials
	client, err := google.DefaultClient(
		oauth2.NoContext,
		"https://www.googleapis.com/auth/devstorage.read_only")
	if err != nil {
		return nil, err
	}

	// Create the Google Cloud Storage service
	service, err := storage.New(client)
	if err != nil {
		return nil, err
	}

	// Create the request to list buckets for the project id
	request := service.Buckets.List(projectId)

	// Execute the request
	buckets, err := request.Do()
	if err != nil {
		return nil, err
	}

	return buckets.Items, nil
}

// [END ListBuckets]

// main will simply retrieve a list of buckets and print them.
func main() {
	buckets, err := ListBuckets(os.Getenv("TEST_PROJECT_ID"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Print out the results
	for _, bucket := range buckets {
		fmt.Println(bucket.Name)
	}
}
