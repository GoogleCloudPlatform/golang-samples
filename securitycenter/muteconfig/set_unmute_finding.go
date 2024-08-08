// Copyright 2023 Google LLC
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

package muteconfig

// [START securitycenter_set_unmute]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
)

// setMute mutes an individual finding, can also unmute or reset the mute state of a finding.
// If a finding is already muted, muting it again has no effect.
// Various mute states are: UNDEFINED/MUTE/UNMUTE.
func setUnmute(w io.Writer, findingPath string) error {
	// findingPath: The relative resource name of the finding. See:
	// https://cloud.google.com/apis/design/resource_names#relative_resource_name
	// Use any one of the following formats:
	//  - organizations/{organization_id}/sources/{source_id}/finding/{finding_id}
	//  - folders/{folder_id}/sources/{source_id}/finding/{finding_id}
	//  - projects/{project_id}/sources/{source_id}/finding/{finding_id}
	// findingPath := fmt.Sprintf("projects/%s/sources/%s/finding/%s", "your-google-cloud-project-id", "source", "finding-id")
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.SetMuteRequest{
		Name: findingPath,
		Mute: securitycenterpb.Finding_UNMUTED}

	finding, err := client.SetMute(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to set the specified mute value: %w", err)
	}
	fmt.Fprintf(w, "Mute value for the finding: %s is %s", finding.Name, finding.Mute)
	return nil
}

// [END securitycenter_set_unmute]
