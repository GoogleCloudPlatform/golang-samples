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

// [START kms_get_retired_resource]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
)

// getRetiredResource gets a retired resource.
func getRetiredResource(w io.Writer, name string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"
	//
	// You can also get a retired resource by its retired resource name:
	// name := "projects/my-project/locations/us-east1/retiredResources/my-retired-resource"

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.GetRetiredResourceRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetRetiredResource(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get retired resource: %w", err)
	}

	fmt.Fprintf(w, "Got retired resource: %s\n", result.Name)
	return nil
}

// [END kms_get_retired_resource]
