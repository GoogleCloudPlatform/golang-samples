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

// [START asset_quickstart_list_saved_queries]

package list

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"google.golang.org/api/iterator"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func listSavedQueries(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("asset.NewClient: %w", err)
	}
	defer client.Close()

	cloudResourceManagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return fmt.Errorf("cloudresourcemanager.NewService: %w", err)
	}

	project, err := cloudResourceManagerClient.Projects.Get(projectID).Do()
	if err != nil {
		return fmt.Errorf("cloudResourceManagerClient.Projects.Get.Do: %w", err)
	}
	projectNumber := strconv.FormatInt(project.ProjectNumber, 10)
	// query name is defined as 'projects/PROJECT_NUMBER'/savedQueries/SAVED_QUERY_ID.
	// we should translate the projectId into a project number first.
	parent := fmt.Sprintf("projects/%s", projectNumber)
	req := &assetpb.ListSavedQueriesRequest{
		Parent: parent}
	it := client.ListSavedQueries(ctx, req)
	for {
		response, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error getting saved queries:%w", err)
		}
		fmt.Fprintf(w, "Query Name: %s\n", response.Name)
		fmt.Fprintf(w, "Query Description:%s\n", response.Description)
		fmt.Fprintf(w, "Query Content:%s\n", response.Content)
	}
	return nil
}

// [END asset_quickstart_list_saved_queries]
