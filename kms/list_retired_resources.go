// Copyright 2026 Google LLC
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

package kms

// [START kms_list_retired_resources]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
	"google.golang.org/api/iterator"
)

// listRetiredResources lists retired resources.
func listRetiredResources(w io.Writer, parent string) error {
	// parent := "projects/my-project/locations/us-east1"

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.ListRetiredResourcesRequest{
		Parent: parent,
	}

	// Call the API.
	it := client.ListRetiredResources(ctx, req)

	// Iterate over the results.
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list retired resources: %w", err)
		}

		fmt.Fprintf(w, "Retired resource: %s\n", resp.Name)
	}
	return nil
}

// [END kms_list_retired_resources]
