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
 
package main
 
import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
 
	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)
 
func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	env := map[string]string{"GOOGLE_CLOUD_PROJECT": tc.ProjectID}
	queryId := fmt.Sprintf("query-%s", strconv.FormatInt(time.Now().UnixNano(), 10))
 
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		t.Fatalf("asset.NewClient: %v", err)
	}
 
	cloudresourcemanagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		t.Fatalf("cloudresourcemanager.NewService: %v", err)
	}
 
	project, err := cloudresourcemanagerClient.Projects.Get(tc.ProjectID).Do()
	if err != nil {
		t.Fatalf("cloudresourcemanagerClient.Projects.Get.Do: %v", err)
	}
	projectNumber := strconv.FormatInt(project.ProjectNumber, 10)
	parent := fmt.Sprintf("projects/%s", tc.ProjectID)
 
 
	req := &assetpb.CreateSavedQueryRequest{
		Parent: parent,
		SavedQueryId: queryId,
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
 
	m := testutil.BuildMain(t)
	defer m.Cleanup()
 
	if !m.Built() {
		t.Errorf("failed to build app")
	}
 
	stdOut, stdErr, err := m.Run(env, 2*time.Minute, fmt.Sprintf("--saved_query_id=%s", queryId))
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}
	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
	got := string(stdOut)
	if !strings.Contains(got, queryId) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, queryId)
	}
 
	client.DeleteSavedQuery(ctx, &assetpb.DeleteSavedQueryRequest{
		Name: fmt.Sprintf("projects/%s/savedQueries/%s", projectNumber, queryId),
	})
}