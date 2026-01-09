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

package regional_secretmanager

import (
	"context"
	"fmt"
	"io"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// [START secretmanager_list_regional_secret_tag_bindings]

// listRegionalSecretTagBindings lists tag bindings for a regional secret.
func listRegionalSecretTagBindings(w io.Writer, secretName, locationID string) error {
	// secretName := "projects/my-project/locations/us-central1/secrets/my-secret"
	// locationID := "us-central1"

	ctx := context.Background()
	rmEndpoint := fmt.Sprintf("%s-cloudresourcemanager.googleapis.com:443", locationID)
	tagBindingsClient, err := resourcemanager.NewTagBindingsClient(ctx, option.WithEndpoint(rmEndpoint))
	if err != nil {
		return fmt.Errorf("failed to create tagbindings client: %w", err)
	}
	defer tagBindingsClient.Close()

	parent := "//secretmanager.googleapis.com/" + secretName

	it := tagBindingsClient.ListTagBindings(ctx, &resourcemanagerpb.ListTagBindingsRequest{
		Parent: parent,
	})

	fmt.Fprintf(w, "Tag bindings for %s:\n", secretName)
	count := 0
	for {
		binding, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate tag bindings: %w", err)
		}
		fmt.Fprintf(w, "- Tag Value: %s\n", binding.GetTagValue())
		count++
	}
	if count == 0 {
		fmt.Fprintf(w, "No tag bindings found for %s.\n", secretName)
	}

	return nil
}

// [END secretmanager_list_regional_secret_tag_bindings]
