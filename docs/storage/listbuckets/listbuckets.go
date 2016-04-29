// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command listbuckets lists the Google Cloud buckets for a given project ID.
package main

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/api/storage/v1"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

// ListBuckets returns a slice of all the buckets for a given project.
func ListBuckets(projectID string) ([]*storage.Bucket, error) {
	ctx := context.Background()

	// Create the client that uses Application Default Credentials
	// See https://developers.google.com/identity/protocols/application-default-credentials
	client, err := google.DefaultClient(ctx, storage.DevstorageReadOnlyScope)
	if err != nil {
		return nil, err
	}

	// Create the Google Cloud Storage service
	service, err := storage.New(client)
	if err != nil {
		return nil, err
	}

	buckets, err := service.Buckets.List(projectID).Do()
	if err != nil {
		return nil, err
	}

	return buckets.Items, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: listbuckets <projectID>")
		os.Exit(1)
	}
	project := os.Args[1]

	buckets, err := ListBuckets(project)
	if err != nil {
		log.Fatal(err)
	}

	// Print out the results
	for _, bucket := range buckets {
		fmt.Println(bucket.Name)
	}
}
