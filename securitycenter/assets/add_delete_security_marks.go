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

// Pacakge assets contains example snippets for listing CSCC Assets.
package assets

// [START add_delete_security_marks]
import (
	"context"
	"fmt"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
	"google.golang.org/genproto/protobuf/field_mask"
)

// addDeleteSecurityMarks adds/updates "key_a" and deletes  "key_b" from
// assetName's securityMarks, and returns the updated marks.  assetName is the
// resource path for an asset.
func addDeleteSecurityMarks(assetName string) (*securitycenterpb.SecurityMarks, error) {
	// assetName := "organizations/123123342/assets/12312321"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error instantiating client %v\n", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	req := &securitycenterpb.UpdateSecurityMarksRequest{
		// If not set or empty, all marks would be cleared.
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"marks.key_a", "marks.key_b"},
		},
		SecurityMarks: &securitycenterpb.SecurityMarks{
			Name: fmt.Sprintf("%s/securityMarks", assetName),
			// Intentionally not setting marks for key_b to
			// delete it.
			Marks: map[string]string{"key_a": "new_value_a"},
		},
	}
	return client.UpdateSecurityMarks(ctx, req)
}

// [END add_delete_security_marks]
