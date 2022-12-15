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

// [START asset_quickstart_get_saved_query]

package get

import (
	"context"
	"fmt"
	"io"
	"strconv"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func getSavedQuery(w io.Writer, projectId, savedQueryID string) error {
	// projectID := "my-project-id"
	// savedQueryID := "query-ID"
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("asset.NewClient: %v", err)
	}
	defer client.Close()

	cloudResourceManagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return fmt.Errorf("cloudresourcemanager.NewService: %v", err)
	}

	project, err := cloudResourceManagerClient.Projects.Get(projectId).Do()
	if err != nil {
		return fmt.Errorf("cloudresourcemanagerClient.Projects.Get.Do: %v", err)
	}
	projectNumber := strconv.FormatInt(project.ProjectNumber, 10)
	// query name is defined as 'projects/PROJECT_NUMBER'/savedQueries/SAVED_QUERY_ID.
	// hence we should translate the projectId into a project number first.
	queryName := fmt.Sprintf("projects/%s/savedQueries/%s", projectNumber, savedQueryID)
	req := &assetpb.GetSavedQueryRequest{
		Name: queryName}
	response, err := client.GetSavedQuery(ctx, req)
	if err != nil {
		return fmt.Errorf("client.GetSavedQuery: %v", err)
	}
	fmt.Fprintf(w, "Query Name: %s\n", response.Name)
	fmt.Fprintf(w, "Query Description:%s\n", response.Description)
	fmt.Fprintf(w, "Query Content:%s\n", response.Content)
	return nil
}

// [END asset_quickstart_get_saved_query]
