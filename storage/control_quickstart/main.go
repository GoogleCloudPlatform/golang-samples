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

// [START storage_control_quickstart_sample]
// This sample demonstrates how to set up and make an API call with the
// Storage Control client.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	controlpb "cloud.google.com/go/storage/control/apiv2/controlpb"
)

func main() {
	// Set this flag to an existing Cloud Storage bucket when running the sample.
	bucketName := flag.String("bucket", "", "Cloud Storage bucket name")
	flag.Parse()
	log.Printf("bucket : %v", *bucketName)

	ctx := context.Background()

	// Create a client.
	client, err := control.NewStorageControlClient(ctx)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Create a request to get the storage layout for the bucket.
	req := &controlpb.GetStorageLayoutRequest{
		Name: fmt.Sprintf("projects/_/buckets/%v/storageLayout", *bucketName),
	}

	// Set a context timeout and send the request.
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := client.GetStorageLayout(ctx, req)
	if err != nil {
		log.Fatalf("GetStorageLayout: %v", err)
	}

	// Use response.
	fmt.Printf("Bucket %v has location type %v", bucketName, res.LocationType)
}

// [END storage_control_quickstart_sample]
