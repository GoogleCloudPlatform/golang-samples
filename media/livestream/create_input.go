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

// [START livestream_create_input]
import (
	"context"
	"fmt"
	"io"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
)

// createInput creates an input endpoint. You send an input video stream to this
// endpoint.
func createInput(w io.Writer, projectID, location, inputID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputID := "my-input"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.CreateInputRequest{
		Parent:  fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		InputId: inputID,
		Input: &livestreampb.Input{
			Type: livestreampb.Input_RTMP_PUSH,
		},
	}
	// Creates the input.
	op, err := client.CreateInput(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateInput: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Input: %v", response.Name)
	return nil
}

// [END livestream_create_input]
