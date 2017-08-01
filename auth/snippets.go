// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contain Google Cloud authentication snippets.
package snippets

import (
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func adc() {
	ctx := context.Background()

	// [START auth_cloud_implicit]
	// Instantiate a client. If you don't specify credentials
	// when constructing the client, the client library will look
	// for credentials in the environment.
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	it := storageClient.Buckets(ctx, "project-id")
	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(bucketAttrs.Name)
	}
	// [END auth_cloud_implicit]
}
