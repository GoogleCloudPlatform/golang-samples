// Copyright 2022 Google LLC
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

// [START asset_quickstart_create_saved_query]

// Sample create-saved-query create saved query.
package create

import (
	"context"
	"fmt"
	"io"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
)

func createSavedQuery(w io.Writer, projectID, savedQueryID string) error {
	// projectID := "my-project-id"
	// savedQueryID := "query-ID"
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("asset.NewClient: %w", err)
	}
	defer client.Close()
	parent := fmt.Sprintf("projects/%s", projectID)
	req := &assetpb.CreateSavedQueryRequest{
		Parent:       parent,
		SavedQueryId: savedQueryID,
		SavedQuery: &assetpb.SavedQuery{
			Content: &assetpb.SavedQuery_QueryContent{
				QueryContent: &assetpb.SavedQuery_QueryContent_IamPolicyAnalysisQuery{
					IamPolicyAnalysisQuery: &assetpb.IamPolicyAnalysisQuery{
						Scope: parent,
						AccessSelector: &assetpb.IamPolicyAnalysisQuery_AccessSelector{
							Permissions: []string{"iam.serviceAccount.actAs"},
						},
					},
				},
			},
		}}
	response, err := client.CreateSavedQuery(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateSavedQuery: %w", err)
	}
	fmt.Fprintf(w, "Query Name: %s\n", response.Name)
	fmt.Fprintf(w, "Query Description:%s\n", response.Description)
	fmt.Fprintf(w, "Query Content:%s\n", response.Content)
	return nil
}

// [END asset_quickstart_create_saved_query]
