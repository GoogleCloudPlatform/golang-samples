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

// [START asset_quickstart_delete_saved_query]

package delete

import (
	"context"
	"fmt"
	"io"
	"strconv"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func deleteSavedQuery(w io.Writer, projectID, savedQueryID string) error {
	// projectID := "my-project-id"
	// savedQueryID := "query-ID"
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
	savedQueryName := fmt.Sprintf("projects/%s/savedQueries/%s", projectNumber, savedQueryID)
	req := &assetpb.DeleteSavedQueryRequest{
		Name: savedQueryName,
	}
	if err = client.DeleteSavedQuery(ctx, req); err != nil {
		return fmt.Errorf("client.DeleteSavedQuery: %w", err)
	}
	fmt.Fprintf(w, "Deleted Saved Query")
	return nil
}

// [END asset_quickstart_delete_saved_query]
