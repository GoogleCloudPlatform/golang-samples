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

// [START livestream_update_pool]
import (
	"context"
	"fmt"
	"io"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updatePool updates the pool's peered network.
func updatePool(w io.Writer, projectID, location, poolID, peeredNetwork string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// poolID := "default"
	// peeredNetwork :=  "projects/my-network-project-number/global/networks/my-network-name"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.UpdatePoolRequest{
		Pool: &livestreampb.Pool{
			Name: fmt.Sprintf("projects/%s/locations/%s/pools/%s", projectID, location, poolID),
			NetworkConfig: &livestreampb.Pool_NetworkConfig{
				PeeredNetwork: peeredNetwork,
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"network_config",
			},
		},
	}
	// Updates the pool.
	op, err := client.UpdatePool(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdatePool: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Updated pool: %v", response.Name)
	return nil
}

// [END livestream_update_pool]
