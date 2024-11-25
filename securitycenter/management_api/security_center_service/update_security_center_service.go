// Copyright 2024 Google LLC
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

package security_center_service

import (
	"context"
	"fmt"
	"io"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateSecurityCenterService updates a Security Center service configuration.
func updateSecurityCenterService(w io.Writer, parent string, serviceID string) error {
	// parent: Use any one of the following options:
	// - organizations/{organization_id}/locations/{location_id}
	// - folders/{folder_id}/locations/{location_id}
	// - projects/{project_id}/locations/{location_id}
	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycentermanagement.NewClient: %w", err)
	}
	defer client.Close()

	// Prepare the updated Security Center service object.
	service := &securitycentermanagementpb.SecurityCenterService{
		Name:                    fmt.Sprintf("%s/securityCenterServices/%s", parent, serviceID),
		IntendedEnablementState: securitycentermanagementpb.SecurityCenterService_ENABLED,
	}

	// Specify which fields to update using a FieldMask.
	updateMask := &fieldmaskpb.FieldMask{
		Paths: []string{"intended_enablement_state"},
	}

	// Create the update request.
	req := &securitycentermanagementpb.UpdateSecurityCenterServiceRequest{
		SecurityCenterService: service,
		UpdateMask:            updateMask,
	}

	// Execute the update request.
	updatedService, err := client.UpdateSecurityCenterService(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update SecurityCenterService: %w", err)
	}

	fmt.Fprintf(w, "Updated SecurityCenterService: %s with new enablement state: %v\n", updatedService.Name, updatedService.IntendedEnablementState)
	return nil
}
