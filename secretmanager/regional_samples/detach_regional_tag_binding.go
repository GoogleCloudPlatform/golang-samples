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
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// [START secretmanager_detach_regional_tag_binding]

// detachRegionalTag detaches a tag value from a regional secret.
func detachRegionalTag(w io.Writer, secretName, locationID, tagValue string) error {
	// secretName := "projects/my-project/locations/us-central1/secrets/my-secret"
	// locationID := "us-central1"
	// tagValue := "tagValues/123456789012"

	ctx := context.Background()
	rmEndpoint := fmt.Sprintf("%s-cloudresourcemanager.googleapis.com:443", locationID)
	client, err := resourcemanager.NewTagBindingsClient(ctx, option.WithEndpoint(rmEndpoint))
	if err != nil {
		return fmt.Errorf("failed to create tagbindings client: %w", err)
	}
	defer client.Close()

	parent := "//secretmanager.googleapis.com/" + secretName

	var bindingName string
	it := client.ListTagBindings(ctx, &resourcemanagerpb.ListTagBindingsRequest{
		Parent: parent,
	})
	for {
		binding, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate tag bindings: %w", err)
		}
		if binding.GetTagValue() == tagValue {
			bindingName = binding.GetName()
			break
		}
	}

	if bindingName == "" {
		fmt.Fprintf(w, "Tag binding for value %s not found on %s.\n", tagValue, secretName)
		return nil
	}

	op, err := client.DeleteTagBinding(ctx, &resourcemanagerpb.DeleteTagBindingRequest{
		Name: bindingName,
	})
	if err != nil {
		return fmt.Errorf("failed to detach tag binding: %w", err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("failed to wait for tag binding deletion: %w", err)
	}

	fmt.Fprintf(w, "Detached tag value %s from %s\n", tagValue, secretName)
	return nil
}

// [END secretmanager_detach_regional_tag_binding]
