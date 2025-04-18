// Copyright 2025 Google LLC
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

// Sample code for deleting a model armor template.

package modelarmor

// [START modelarmor_delete_template]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// deleteModelArmorTemplate method deletes a Model Armor template with the provided ID.
//
// w io.Writer: The writer to use for logging.
// projectID string: The ID of the Google Cloud project.
// locationID string: The ID of the Google Cloud location.
// templateID string: The ID of the template to delete.
func deleteModelArmorTemplate(w io.Writer, projectID, locationID, templateID string) error {
	ctx := context.Background()

	// Create option for Model Armor client.
	opts := option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID))
	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create client for project %s, location %s: %w", projectID, locationID, err)
	}
	defer client.Close()

	// Build the request for deleting the template.
	req := &modelarmorpb.DeleteTemplateRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID),
	}

	// Delete the template.
	if err := client.DeleteTemplate(ctx, req); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	// Print the success message using fmt.Fprintf with the io.Writer.
	fmt.Fprintf(w, "Successfully deleted Model Armor template: %s\n", req.Name)

	return err
}

// [END modelarmor_delete_template]
