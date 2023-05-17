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

// [START data_catalog_ptm_create_taxonomy]
import (
	"context"
	"fmt"
	"io"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"cloud.google.com/go/datacatalog/apiv1beta1/datacatalogpb"
)

// createTaxonomy creates a taxonomy resource, which holds the collection of policy tags.
//
// This example marks the taxonomy as valid for defining fine grained access control, also
// known as column-level access control when used in conjunction with BigQuery.
func createTaxonomy(w io.Writer, projectID, location, displayName string) (string, error) {
	// projectID := "my-project-id"
	// location := "us"
	// displayName := "example-taxonomy"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return "", fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %w", err)
	}
	defer policyClient.Close()

	req := &datacatalogpb.CreateTaxonomyRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Taxonomy: &datacatalogpb.Taxonomy{
			DisplayName: displayName,
			Description: "Taxonomy created via basic snippet testing",
			ActivatedPolicyTypes: []datacatalogpb.Taxonomy_PolicyType{
				datacatalogpb.Taxonomy_FINE_GRAINED_ACCESS_CONTROL,
			},
		},
	}
	resp, err := policyClient.CreateTaxonomy(ctx, req)
	if err != nil {
		return "", fmt.Errorf("CreateTaxonomy: %w", err)
	}

	fmt.Fprintf(w, "Taxonomy %s was created.\n", resp.Name)
	return resp.Name, nil
}

// [END data_catalog_ptm_create_taxonomy]
