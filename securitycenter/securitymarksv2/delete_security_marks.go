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

package securitymarksv2

// [START securitycenter_delete_security_marks_v2]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// deleteSecurityMarks deletes security marks "key_a" and  "key_b" from
// assetName's marks. assetName is the resource path for an asset.
func deleteSecurityMarks(w io.Writer, assetName string) error {
	// Specify the value of 'assetName' in one of the following formats:
	// 		assetName := "organizations/{org_id}/assets/{asset_id}"
	//		assetName := "projects/{project_id}/assets/{asset_id}"
	//		assetName := "folders/{folder_id}/assets/{asset_id}"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	req := &securitycenterpb.UpdateSecurityMarksRequest{
		// If not set or empty, all marks would be cleared.
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"marks.key_a", "marks.key_b"},
		},
		SecurityMarks: &securitycenterpb.SecurityMarks{
			Name: fmt.Sprintf("%s/securityMarks", assetName),
			// Intentionally not setting marks with the
			// corresponding field mask deletes them.
		},
	}
	updatedMarks, err := client.UpdateSecurityMarks(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateSecurityMarks: %w", err)
	}

	fmt.Fprintf(w, "Updated marks: %s\n", updatedMarks.Name)
	for k, v := range updatedMarks.Marks {
		fmt.Fprintf(w, "%s = %s\n", k, v)
	}
	return nil
}

// [END securitycenter_delete_security_marks_v2]
