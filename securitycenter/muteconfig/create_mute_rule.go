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

// [START securitycenter_create_mute_config]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
)

// createMuteRule: Creates a mute configuration under a given scope that will mute
// all new findings that match a given filter.
// Existing findings will not be muted.
func createMuteRule(w io.Writer, parent string, muteConfigId string) error {
	// parent: Use any one of the following options:
	//             - organizations/{organization_id}
	//             - folders/{folder_id}
	//             - projects/{project_id}
	// parent := fmt.Sprintf("projects/%s", "your-google-cloud-project-id")
	// muteConfigId: Set a random id; max of 63 chars.
	// muteConfigId := "random-mute-id-" + uuid.New().String()
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	muteConfig := &securitycenterpb.MuteConfig{
		Description: "Mute low-medium IAM grants excluding 'compute' ",
		// Set mute rule(s).
		// To construct mute rules and for supported properties, see:
		// https://cloud.google.com/security-command-center/docs/how-to-mute-findings#create_mute_rules
		Filter: "severity=\"LOW\" OR severity=\"MEDIUM\" AND " +
			"category=\"Persistence: IAM Anomalous Grant\" AND " +
			"-resource.type:\"compute\"",
	}

	req := &securitycenterpb.CreateMuteConfigRequest{
		Parent:       parent,
		MuteConfigId: muteConfigId,
		MuteConfig:   muteConfig,
	}

	response, err := client.CreateMuteConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create mute rule: %w", err)
	}
	fmt.Fprintf(w, "Mute rule created successfully: %s", response.Name)
	return nil
}

// [END securitycenter_create_mute_config]
