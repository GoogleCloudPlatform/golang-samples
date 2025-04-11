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

// [START data_catalog_ptm_delete_policytag]
import (
	"context"
	"fmt"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"cloud.google.com/go/datacatalog/apiv1beta1/datacatalogpb"
)

// deletePolicyTag deletes a given policy tag.
func deletePolicyTag(policyTagID string) error {
	// policyTagID := "projects/myproject/locations/us/taxonomies/1234/policyTags/5678"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %w", err)
	}
	defer policyClient.Close()

	req := &datacatalogpb.DeletePolicyTagRequest{
		Name: policyTagID,
	}
	return policyClient.DeletePolicyTag(ctx, req)
}

// [END data_catalog_ptm_delete_policytag]
