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

package findingsv2

// [START securitycenter_add_finding_security_marks_v2]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// addSecurityMarks adds/updates security marks for the findingName and
// returns the updated marks. Specifically, it sets "key_a" an "key_b" to
// "value_a" and "value_b" respectively. findingName is the resource path for
// the finding to add marks to.
func addSecurityMarks(w io.Writer, findingName string) error {
	// Specify the value of 'findingName' in one of the following formats:
	// 		"organizations/{orgId}/sources/{sourceId}/findings/{findingId}"
	// 		"projects/{projectId}/sources/{sourceId}/findings/{findingId}"
	// 		"folders/{folderId}/sources/{sourceId}/findings/{findingId}"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	req := &securitycenterpb.UpdateSecurityMarksRequest{
		// If not set or empty, all marks would be cleared before
		// adding the new marks below.
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"marks.key_a", "marks.key_b"},
		},
		SecurityMarks: &securitycenterpb.SecurityMarks{
			Name: fmt.Sprintf("%s/securityMarks", findingName),
			// Note keys correspond to the last part of each path.
			Marks: map[string]string{"key_a": "value_a", "key_b": "value_b"},
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

// [END securitycenter_add_finding_security_marks_v2]
