// Copyright 2026 Google LLC
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

package regional_parametermanager

// [START parametermanager_create_regional_param_template_version]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/option"
)

// createRegionalParamTemplateVersion creates a new version of a parameter template.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// locationID: The region where the resources are located.
// templateID: The ID of the template.
// versionID: The ID of the template version.
// payload: The payload of the version.
//
// The function returns an error if the operation fails.
func createRegionalParamTemplateVersion(w io.Writer, projectID, locationID, templateID, versionID, payload string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the parent template.
	parent := fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID)

	// Build the request to create a new template version.
	req := &parametermanagerpb.CreateTemplateVersionRequest{
		Parent:            parent,
		TemplateVersionId: versionID,
		TemplateVersion: &parametermanagerpb.TemplateVersion{
			Payload: &parametermanagerpb.TemplateVersionPayload{
				Data: []byte(payload),
			},
		},
	}

	// Call the API to create the template version.
	version, err := client.CreateTemplateVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create template version: %w", err)
	}

	fmt.Fprintf(w, "Created regional parameter template version: %s\n", version.Name)
	return nil
}

// [END parametermanager_create_regional_param_template_version]
