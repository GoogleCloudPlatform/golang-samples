// Copyright 2019 Google LLC
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

// settings contains snippets for working with CSCC organization settings.
package settings

// [START get_org_settings]
import (
	"context"
	"fmt"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
	"google.golang.org/genproto/protobuf/field_mask"
)

// Turns on asset discovery for orgID and returns the updated organization
// settings. orgID is the numeric Organization ID.
func enableAssetDiscoery(orgID string) (*securitycenterpb.OrganizationSettings, error) {
	// orgID := "12321311"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error instantiating client %v\n", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	req := &securitycenterpb.UpdateOrganizationSettingsRequest{
		OrganizationSettings: &securitycenterpb.OrganizationSettings{
			Name:                 fmt.Sprintf("organizations/%s/organizationSettings", orgID),
			EnableAssetDiscovery: true,
		},
		// Only update the asset discovery setting.
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"enable_asset_discovery"},
		},
	}
	settings, err := client.UpdateOrganizationSettings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Error creating finding: %v", err)
	}
	return settings, nil
}

// [END get_org_settings]
