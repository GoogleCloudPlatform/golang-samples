// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Prompt: You are an excellent Go programmer.
// Write a Go language program to delete a dataset.

package main

import (
	"context"
	"fmt"
	"log"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1"
)

func main() {
	ctx := context.Background()
	aiplatformService, err := aiplatform.NewDatasetClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	defer aiplatformService.Close()

	req := &aiplatformpb.DeleteDatasetRequest{
		Name: "projects/123456789/locations/us-central1/datasets/123456789",
	}

	op, err := aiplatformService.DeleteDataset(ctx, req)
	if err != nil {
		log.Fatalf("Failed to delete dataset: %v\n", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		log.Fatalf("Failed to wait for operation: %v\n", err)
	}

	fmt.Printf("Deleted dataset: %s\n", resp.GetName())
}

// [END aiplatform_delete_dataset_sample]
