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

// [START data_catalog_ptm_list_taxonomies]
import (
	"context"
	"fmt"
	"io"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"cloud.google.com/go/datacatalog/apiv1beta1/datacatalogpb"
	"google.golang.org/api/iterator"
)

// listTaxonomies prints information about the taxonomies contained within a specific
// project and location.
func listTaxonomies(w io.Writer, projectID, location string) error {
	// projectID := "my-project-id"
	// location := "us"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %w", err)
	}
	defer policyClient.Close()

	req := &datacatalogpb.ListTaxonomiesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}
	it := policyClient.ListTaxonomies(ctx, req)
	fmt.Fprintf(w, "listing taxonomies in project %s and location %s\n", projectID, location)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListTaxonomies iteration error: %w", err)
		}

		fmt.Fprintf(w, "\t- %s: %s\n", resp.Name, resp.DisplayName)
	}
	return nil
}

// [END data_catalog_ptm_list_taxonomies]
