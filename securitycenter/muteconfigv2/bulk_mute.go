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

package muteconfigv2

// [START securitycenter_bulk_mute_v2]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
)

// bulkMute kicks off a long-running operation (LRO) to bulk mute findings for a parent based on a filter.
// The parent can be either an organization, folder, or project. The findings
// matched by the filter will be muted after the LRO is done.
func bulkMute(w io.Writer, parent string, muteRule string) error {
	// parent: Use any one of the following options:
	//             - organizations/{organization_id}
	//             - folders/{folder_id}
	//             - projects/{project_id}
	// parent := fmt.Sprintf("projects/%s", "your-google-cloud-project-id")
	// muteRule: Expression that identifies findings that should be muted.
	// To create mute rules, see:
	// https://cloud.google.com/security-command-center/docs/how-to-mute-findings#create_mute_rules
	// muteRule := "filter-condition"
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.BulkMuteFindingsRequest{
		Parent: parent,
		Filter: muteRule,
	}

	op, err := client.BulkMuteFindings(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to bulk mute findings: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to bulk mute findings: %w", err)
	}
	fmt.Fprintf(w, "Bulk mute findings completed successfully! %s", response)
	return nil
}

// [END securitycenter_bulk_mute_v2]
