// Copyright 2020 Google LLC
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

package policytagmanager

// [START data_catalog_ptm_list_policytags]
import (
	"context"
	"fmt"
	"io"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"cloud.google.com/go/datacatalog/apiv1beta1/datacatalogpb"
	"google.golang.org/api/iterator"
)

// listPolicyTags prints information about the policy tags within a given taxonomy
// resource.
func listPolicyTags(w io.Writer, parentTaxonomyID string) error {
	// parentTaxonomyID := projects/myproject/locations/us/taxonomies/1234"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %w", err)
	}
	defer policyClient.Close()

	req := &datacatalogpb.ListPolicyTagsRequest{
		Parent: parentTaxonomyID,
	}
	it := policyClient.ListPolicyTags(ctx, req)
	fmt.Fprintf(w, "listing policy tags in taxonomy %s\n", parentTaxonomyID)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListPolicyTags iteration error: %w", err)
		}

		fmt.Fprintf(w, "\t- %s (%s)", resp.Name, resp.DisplayName)
		if resp.ParentPolicyTag != "" {
			fmt.Fprintf(w, " has parent Tag %s", resp.ParentPolicyTag)
		}
		fmt.Fprintln(w)
	}
	return nil
}

// [END data_catalog_ptm_list_policytags]
