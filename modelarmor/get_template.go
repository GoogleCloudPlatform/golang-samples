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

// Sample code for getting a model armor template.

package modelarmor

// [START modelarmor_get_template]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// getModelArmorTemplate gets a Model Armor template.
//
// This method retrieves a Model Armor template.
//
// Args:
//
//	w io.Writer: The writer to use for logging.
//	projectID string: The ID of the project.
//	locationID string: The location of the template.
//	templateID string: The ID of the template.
func getModelArmorTemplate(w io.Writer, projectID, locationID, templateID string) error {
	ctx := context.Background()

	// Create options for the Model Armor client/
	opts := option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID))
	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create client for project %s, location %s: %w", projectID, locationID, err)
	}
	defer client.Close()

	// Initialize request arguments.
	req := &modelarmorpb.GetTemplateRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID),
	}

	// Get the template.
	response, err := client.GetTemplate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Print the template name using fmt.Fprintf with the io.Writer.
	fmt.Fprintf(w, "Retrieved template: %s\n", response.Name)

	return nil
}

// [END modelarmor_get_template]
