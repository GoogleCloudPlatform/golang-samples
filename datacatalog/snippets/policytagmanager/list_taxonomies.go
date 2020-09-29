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

import (
	"context"
	"fmt"
	"io"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	datacatalogpb "google.golang.org/genproto/googleapis/cloud/datacatalog/v1beta1"
)

func listTaxonomies(projectID string, w io.Writer) error {
	// projectID := "my-project-id"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx,
		option.WithEndpoint("staging-datacatalog.sandbox.googleapis.com:443"))
	if err != nil {
		return fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %v", err)
	}
	defer policyClient.Close()

	req := &datacatalogpb.ListTaxonomiesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/us", projectID),
	}
	it := policyClient.ListTaxonomies(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListTaxonomies iteration error: %v", err)
		}

		fmt.Fprintf(w, "taxonomy %s: %s\n", resp.Name, resp.DisplayName)
	}
	return nil
}
