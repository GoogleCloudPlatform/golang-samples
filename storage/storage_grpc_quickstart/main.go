// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START storage_grpc_quickstart]

// Sample storage-quickstart creates a Google Cloud Storage bucket using
// gRPC API.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
)

func main() {
	ctx := context.Background()

	// Expects Google Cloud Platform project ID and Cloud Storage bucket
	projectID := flag.String("project", "", "Cloud Platform project id")
	bucketName := flag.String("bucket", "", "Cloud Storage bucket name")
	flag.Parse()

	// Creates a gRPC enabled client.
	client, err := storage.NewGRPCClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Creates a Bucket instance.
	bucket := client.Bucket(*bucketName)

	// Creates the new bucket.
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := bucket.Create(ctx, *projectID, nil); err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	fmt.Printf("Bucket %v created.\n", *bucketName)
}

// [END storage_grpc_quickstart]
