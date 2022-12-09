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
 
package CreateSavedQueryFunction
 
import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"os"
 
	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	"github.com/gofrs/uuid"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)
 
func TestMain(t *testing.T) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	savedQueryID := fmt.Sprintf("query-%s", uuid.Must(uuid.NewV4()).String()[:8])
 
	err := createSavedQuery(projectID, savedQueryID);
 
	ctx := context.Background()
	cloudResourceManagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		t.Fatalf("cloudresourcemanager.NewService: %v", err)
	}
 
	project, err := cloudResourceManagerClient.Projects.Get(projectID).Do()
	if err != nil {
		t.Fatalf("cloudresourcemanager.Projects.Get.Do: %v", err)
	}
	// query name is defined as 'projects/PROJECT_NUMBER'/savedQueries/SAVED_QUERY_ID.
	// hence we should translate the projectId into a project number first.
	projectNumber := strconv.FormatInt(project.ProjectNumber, 10)
 
	client, err := asset.NewClient(ctx)
	if err != nil {
		t.Fatalf("asset.NewClient: %v", err)
	}
 
 
	err = client.DeleteSavedQuery(ctx, &assetpb.DeleteSavedQueryRequest{
		Name: fmt.Sprintf("projects/%s/savedQueries/%s", projectNumber, savedQueryID),
	})
 
	if err != nil {
		t.Fatalf("client.DeleteSavedQuery: %v", err);
	}
}