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

package livestream

// [START livestream_update_input]
import (
	"context"
	"fmt"
	"io"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateInput updates an existing input endpoint. This sample adds a
// preprocessing configuration to an existing input.
func updateInput(w io.Writer, projectID, location, inputID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputID := "my-input"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.UpdateInputRequest{
		Input: &livestreampb.Input{
			Name: fmt.Sprintf("projects/%s/locations/%s/inputs/%s", projectID, location, inputID),
			PreprocessingConfig: &livestreampb.PreprocessingConfig{
				Crop: &livestreampb.PreprocessingConfig_Crop{
					TopPixels:    5,
					BottomPixels: 5,
				},
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"preprocessing_config",
			},
		},
	}
	// Updates the input.
	op, err := client.UpdateInput(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateInput: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Updated input: %v", response.Name)
	return nil
}

// [END livestream_update_input]
