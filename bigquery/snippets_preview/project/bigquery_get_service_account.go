// Copyright 2025 Google LLC
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

package project

// [START bigquery_get_service_account]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"
)

// getServiceAccount demonstrates how to interrogate a project to get the service account
// identity for the BigQuery project.  It's often used in conjunction with features related
// to the Google Cloud KMS service and configuration of encryption settings.
func getServiceAccount(client *apiv2_client.Client, w io.Writer, projectID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	ctx := context.Background()

	// Construct a request.
	req := &bigquerypb.GetServiceAccountRequest{
		ProjectId: projectID,
	}
	resp, err := client.GetServiceAccount(ctx, req)
	if err != nil {
		return fmt.Errorf("GetServiceAccount: %w", err)
	}
	// Print the email address of the account to the provided writer.
	fmt.Fprintf(w, "GetServiceAccount email is %q", resp.GetEmail())
	return nil
}

// [END bigquery_get_service_account]
