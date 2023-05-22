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

// [START data_catalog_ptm_get_policytag]
import (
	"context"
	"fmt"
	"io"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"cloud.google.com/go/datacatalog/apiv1beta1/datacatalogpb"
)

// getPolicyTag prints information about a given policy tag.
func getPolicyTag(w io.Writer, policyTagID string) error {
	// policyTagID := "projects/myproject/locations/us/taxonomies/1234/policyTags/5678"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %w", err)
	}
	defer policyClient.Close()

	req := &datacatalogpb.GetPolicyTagRequest{
		Name: policyTagID,
	}
	resp, err := policyClient.GetPolicyTag(ctx, req)
	if err != nil {
		return fmt.Errorf("GetPolicyTag: %w", err)
	}
	fmt.Fprintf(w, "PolicyTag %s has Display Name %s", resp.Name, resp.DisplayName)
	if resp.ParentPolicyTag != "" {
		fmt.Fprintf(w, " and is a child of Policy Tag %s", resp.ParentPolicyTag)
	}
	if len(resp.ChildPolicyTags) > 0 {
		fmt.Fprintf(w, ", with %d child tags", len(resp.ChildPolicyTags))
	}
	fmt.Fprintln(w)
	return nil
}

// [END data_catalog_ptm_get_policytag]
