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

package secretmanager

import (
	"context"
	"fmt"
	"io"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"google.golang.org/api/iterator"
)

// [START secretmanager_list_tag_bindings]

// listTagBindings lists tag bindings attached to a secret.
func listTagBindings(w io.Writer, secretName string) error {
	// secretName := "projects/my-project/secrets/my-secret"

	ctx := context.Background()
	client, err := resourcemanager.NewTagBindingsClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create tagbindings client: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("//secretmanager.googleapis.com/%s", secretName)
	req := &resourcemanagerpb.ListTagBindingsRequest{
		Parent: parent,
	}

	it := client.ListTagBindings(ctx, req)
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

// [END secretmanager_list_tag_bindings]
