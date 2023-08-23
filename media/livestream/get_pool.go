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

// [START livestream_get_pool]
import (
	"context"
	"fmt"
	"io"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
)

// getPool gets a pool.
func getPool(w io.Writer, projectID, location, poolID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// poolID := "default"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.GetPoolRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/pools/%s", projectID, location, poolID),
	}

	response, err := client.GetPool(ctx, req)
	if err != nil {
		return fmt.Errorf("GetPool: %w", err)
	}

	fmt.Fprintf(w, "Pool: %v", response.Name)
	return nil
}

// [END livestream_get_pool]
