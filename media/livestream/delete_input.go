// Copyright 2020 Google LLC
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

// [START livestream_delete_input]
import (
	"context"
	"fmt"
	"io"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
)

// deleteInput deletes a previously-created input endpoint.
func deleteInput(w io.Writer, projectID, location, inputID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputID := "my-input"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.DeleteInputRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/inputs/%s", projectID, location, inputID),
	}

	op, err := client.DeleteInput(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteInput: %w", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Deleted input")
	return nil
}

// [END livestream_delete_input]
