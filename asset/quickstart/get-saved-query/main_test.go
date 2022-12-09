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
 
package GetSavedQueryFunction
 
import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
	"os"
 
	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)
 
func TestGetSavedQuery(t *testing.T) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	savedQueryID := fmt.Sprintf("query-%s", strconv.FormatInt(time.Now().UnixNano(), 10))
	
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		t.Fatalf("asset.NewClient: %v", err)
	}
 
	cloudResourceManagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		t.Fatalf("cloudresourcemanager.NewService: %v", err)
	}
	
	project, err := cloudResourceManagerClient.Projects.Get(projectID).Do()
	if err != nil {
		t.Fatalf("cloudResourceManagerClient.Projects.Get.Do: %v", err)
	}
	// query name is defined as 'projects/PROJECT_NUMBER'/savedQueries/SAVED_QUERY_ID.
	// hence we should translate the projectId into a project number first.
	projectNumber := strconv.FormatInt(project.ProjectNumber, 10)
	parent := fmt.Sprintf("projects/%s", projectID)
 
 
	req := &assetpb.CreateSavedQueryRequest{
		Parent: parent,
		SavedQueryId: savedQueryID,
		SavedQuery: &assetpb.SavedQuery{
			Content: &assetpb.SavedQuery_QueryContent{
				QueryContent: &assetpb.SavedQuery_QueryContent_IamPolicyAnalysisQuery {
					IamPolicyAnalysisQuery: &assetpb.IamPolicyAnalysisQuery {
						Scope: parent,
						AccessSelector: &assetpb.IamPolicyAnalysisQuery_AccessSelector {
							Permissions: []string{"iam.serviceAccount.actAs"},
						},
					},
				},
			},
		}}
	_, err = client.CreateSavedQuery(ctx, req)
	if err != nil {
		t.Fatalf("client.CreateSavedQuery: %v", err)
	}
 
	err = getSavedQuery(projectID, savedQueryID);
 
	client.DeleteSavedQuery(ctx, &assetpb.DeleteSavedQueryRequest{
		Name: fmt.Sprintf("projects/%s/savedQueries/%s", projectNumber, savedQueryID),
	})
}