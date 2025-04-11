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

package livestream

// [START livestream_get_asset]
import (
	"context"
	"fmt"
	"io"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
)

// getAsset gets a previously-created asset.
func getAsset(w io.Writer, projectID, location, assetID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// assetID := "my-asset-id"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.GetAssetRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/assets/%s", projectID, location, assetID),
	}

	response, err := client.GetAsset(ctx, req)
	if err != nil {
		return fmt.Errorf("GetAsset: %w", err)
	}

	fmt.Fprintf(w, "Asset: %v", response.Name)
	return nil
}

// [END livestream_get_asset]
